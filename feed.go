package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

import (
	"appengine"
	"appengine/memcache"
	"appengine/urlfetch"
	"appengine/user"
)

type FeedFormat int

const (
	UNKNOWN FeedFormat = iota
	RSS
	ATOM
	OTHER
)

type FeedInfo struct {
	URL        string
	Subscribed bool
}

type FeedList struct {
	Feeds []FeedInfo
}

type FeedData struct {
	User           string
	Feeds          []FeedCache
	SuggestedFeeds []FeedCache
}

func (FeedData) Template() string { return "feeds.html" }
func (FeedData) Redirect() string { return "" }
func (FeedData) Send() bool       { return true }

type FeedCache struct {
	URL        string
	Title      string
	TimeToLive time.Duration
	Articles   []ArticleCache
}

type SubscriptionCache struct {
	URL    string
	Format FeedFormat
	Update int64
	MD5    string
	Length int
}

type Feed struct {
	URL         string
	Subscribers []string
	Articles    []string
}

type GenericFeed struct {
	XMLName xml.Name
}

func getFeedType(response *http.Response, body []byte) FeedFormat {
	switch response.Header.Get("Content-Type") {
	case "application/atom+xml":
		return ATOM
	case "application/rss+xml":
		return RSS
	default:
		var feed GenericFeed
		xml.Unmarshal(body, &feed)
		switch feed.XMLName.Local {
		case "channel", "rss":
			return RSS
		case "feed":
			return ATOM
		default:
			return OTHER
		}
		return OTHER
	}
	return OTHER
}

const (
	MYDATE    = "1/2/2006 3:04:05 PM"
	ONIONDATE = "Mon, 2 Jan 2006 15:04:05 -0700"
)

func getDate(dateString string) (date time.Time, err error) {
	layouts := []string{time.RFC822, time.RFC822Z, time.RFC3339, time.RFC1123, time.RFC1123Z, time.ANSIC, time.UnixDate, time.RubyDate, ONIONDATE} //, MYDATE
	for _, layout := range layouts {
		date, err = time.Parse(layout, dateString)
		if err == nil && date.Year() != 0 {
			return
		}
	}
	return
}

func getRSS(context appengine.Context, body []byte, url string) (feedCache FeedCache, err error) {
	var rss RSSStruct
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		printError(context, err, url)
	}
	var date time.Time
	for _, channel := range rss.Channel {
		if channel.Ttl > 0 {
			feedCache.TimeToLive = time.Duration(channel.Ttl)
		}
		feedCache.Title = channel.Title
		for _, item := range channel.Item {
			if item.Guid == "" {
				break
			}
			_, err = memcache.Gob.Get(context, item.Guid, nil)
			if err == memcache.ErrCacheMiss {
				date, err = getDate(item.PubDate)
				if err != nil {
					printError(context, fmt.Errorf("rss feed %v has dates that look like %v", channel.Link, item.PubDate), url)
					continue
				}
				var content string
				if item.Content != "" {
					content = item.Content
				} else {
					content = item.Description
				}
				article := ArticleCache{
					URL:     item.Link,
					Title:   item.Title,
					Summary: content,
					ID:      item.Guid,
					Date:    date.Unix(),
					Feed:    feedCache.Title,
				}
				err = memcache.Gob.Set(context, &memcache.Item{Key: item.Guid, Object: article})
				if err != nil {
					printError(context, err, url)
					continue
				}
				feedCache.Articles = append(feedCache.Articles, article)
			} else if err != nil {
				break
			}
		}
	}
	return
}

func getAtom(context appengine.Context, body []byte, url string) (feedCache FeedCache, err error) {
	var feed AtomFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		printError(context, err, url)
	}
	feedCache.Title = feed.Title
	var date time.Time
	for _, item := range feed.Entry {
		if item.Id == "" {
			break
		}
		_, err = memcache.Gob.Get(context, item.Id, nil)
		if err == memcache.ErrCacheMiss {
			date, err = getDate(item.Updated)
			if err != nil {
				printError(context, fmt.Errorf("atom feed %v has dates that look like %v", feed.Link[0].Href, item.Updated), url)
				continue
			}
			var url string
			for _, link := range item.Link {
				if link.Href != "" {
					url = link.Href
					if link.Rel == "alternate" {
						break
					}
				}
			}
			if url == "" {
				break
			}
			article := ArticleCache{
				URL:     url,
				Title:   item.Title,
				Summary: item.Content.Text,
				ID:      item.Id,
				Date:    date.Unix(),
				Feed:    feedCache.Title,
			}
			err = memcache.Gob.Set(context, &memcache.Item{Key: item.Id, Object: article})
			if err != nil {
				printError(context, err, url)
				continue
			}
			feedCache.Articles = append(feedCache.Articles, article)
		} else if err != nil {
			break
		}
	}
	return
}

func getSubscription(context appengine.Context, format FeedFormat, body []byte, url string) (feed FeedCache, err error) { // url is only used for debugging purposes
	switch format {
	case RSS:
		return getRSS(context, body, url)
	case ATOM:
		return getAtom(context, body, url)
	case OTHER:
		err = errors.New(fmt.Sprintf("not a feed")) // later we can delete the feed
	}
	err = errors.New(fmt.Sprintf("could not determine format")) // later we can delete the feed
	return
}

func getSubscriptionURL(context appengine.Context, url string) (feed FeedCache, err error) {
	client := urlfetch.Client(context)
	var response *http.Response
	response, err = client.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()
	var body []byte
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	format := getFeedType(response, body)
	return getSubscription(context, format, body, url)
}

func getSuggestedFeeds(context appengine.Context, userdata UserData) (suggestedFeeds []FeedCache, err error) {
	for _, defaultFeed := range defaultFeeds {
		if !ContainsFeedCache(suggestedFeeds, defaultFeed) && !ContainsString(userdata.Feeds, defaultFeed.URL) {
			suggestedFeeds = append(suggestedFeeds, defaultFeed)
		}
	}
	/*query := datastore.NewQuery("Feed")
	var feed Feed
	for iterator := query.Run(context); ; {
		_, err = iterator.Next(&feed)
		if err == datastore.Done {
			break
		} else if err != nil {
			printError(context, err)
			continue
		}
		if !ContainsString(suggestedFeeds, feed.URL) && !ContainsString(userdata.Feeds, feed.URL) {
			suggestedFeeds = append(suggestedFeeds, feed.URL)
		}
	}*/
	return
}

func feedGET(context appengine.Context, user *user.User, request *http.Request) (data Data, err error) {
	url := request.FormValue("url")
	if url != "" {
		if request.FormValue("unsubscribe") == "1" {
			if user.String() == "default" {
				return
			}
			err = unsubscribe(context, user.String(), url)
			if err != nil {
				return
			}
			var redirect Redirect
			redirect.URL = "/feed"
			return redirect, nil
		} else if request.FormValue("subscribe") == "1" {
			err = subscribeUser(context, user, url)
			if err != nil {
				return
			}
			var redirect Redirect
			redirect.URL = "/feed"
			return redirect, nil
		}
	}
	var feedData FeedData
	if user.String() != "default" {
		feedData.User = user.String()
	}
	var userdata UserData
	_, userdata, err = mustGetUserData(context, user.String())
	if err != nil {
		return
	}

	for _, feed := range userdata.Feeds {
		var item FeedCache
		for _, defaultFeed := range defaultFeeds {
			if defaultFeed.URL == feed {
				item = defaultFeed
				break
			}
		}
		if item.URL == "" {
			_, err = memcache.Gob.Get(context, feed, &item)
			if err == memcache.ErrCacheMiss {
				continue
			}
		}
		feedData.Feeds = append(feedData.Feeds, item)
	}

	feedData.SuggestedFeeds, err = getSuggestedFeeds(context, userdata)
	if err != nil {
		return
	}
	return feedData, nil
}

func feedPOST(context appengine.Context, user *user.User, request *http.Request) (data Data, err error) {
	var body []byte
	body, err = ioutil.ReadAll(request.Body)
	if err != nil {
		return
	}
	var feedList FeedList
	err = json.Unmarshal(body, &feedList)
	if err != nil {
		return
	}
	for _, feed := range feedList.Feeds {
		if feed.Subscribed {
			err = subscribeUser(context, user, feed.URL)
			if err != nil {
				printError(context, err, feed.URL)
				continue
			}
		} else {
			err = unsubscribe(context, user.String(), feed.URL)
			if err != nil {
				printError(context, err, feed.URL)
				continue
			}
		}
	}
	return
}
