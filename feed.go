package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

import (
	"appengine"
	"appengine/datastore"
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

var defaultFeeds = []string{
	"http://www.marco.org/rss",
	"http://www.thegatesnotes.com/RSS",
	"http://www.theverge.com/rss/index.xml",
	"http://daringfireball.net/index.xml",
	"http://scripting.com/rss.xml",
	"http://buzzmachine.com/feed/",
	"http://flowingdata.com/feed/",
	"http://blogs.wsj.com/numbersguy/feed/",
	"http://feeds.kottke.org/main",
	"http://feeds.washingtonpost.com/rss/rss_ezra-klein",
	"http://feeds.feedburner.com/CalculatedRisk",
	"http://feeds.feedburner.com/blogspot/MKuf",
	"http://feeds.feedburner.com/Asymco",
	"http://feeds.feedburner.com/GoogleOperatingSystem",
	"http://feeds.feedburner.com/thebrowser/xrdJ",
	"http://feeds.feedburner.com/538dotcom",
	"http://feeds.feedburner.com/marginalrevolution",
}

type FeedInfo struct {
	URL        string
	Subscribed bool
}

type FeedList struct {
	Feeds []FeedInfo
}

type FeedData struct {
	User           string
	Feeds          []string
	SuggestedFeeds []string
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
	Bytes       []byte   `datastore:",noindex"`
	Subscribers []string `datastore:"-"`
	Articles    []string `datastore:"-"`
}

func (feed *Feed) Load(c <-chan datastore.Property) (err error) {
	for p := range c {
		switch p.Name {
		case "URL":
			feed.URL = p.Value.(string)
		case "Bytes":
			reader := bytes.NewBuffer(p.Value.([]byte))
			decoder := gob.NewDecoder(reader)
			err = decoder.Decode(feed)
			if err != nil {
				return
			}
		}
	}
	return
}

func (feed *Feed) Save(c chan<- datastore.Property) (err error) {
	defer close(c)
	c <- datastore.Property{
		Name:  "URL",
		Value: feed.URL,
	}
	writer := bytes.Buffer{}
	encoder := gob.NewEncoder(&writer)
	err = encoder.Encode(feed)
	if err != nil {
		return
	}
	c <- datastore.Property{Name: "Bytes", Value: writer.Bytes(), NoIndex: true}
	return
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
		err := xml.Unmarshal(body, &feed)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error114: %v\n", err.Error())
		}
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

func getDate(dateString string) (date time.Time, err error) {
	layouts := []string{time.RFC822, time.RFC822Z, time.RFC3339, time.RFC1123, time.RFC1123Z, time.ANSIC, time.UnixDate, time.RubyDate}
	for _, layout := range layouts {
		date, err = time.Parse(layout, dateString)
		if err == nil && date.Year() != 0 {
			return
		}
	}
	return
}

func getRSS(context appengine.Context, body []byte) (feedCache FeedCache, err error) {
	var rss RSSStruct
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error115: %v\n", err.Error())
	}
	err = nil
	var date time.Time
	for _, channel := range rss.Channel {
		if channel.Ttl > 0 {
			feedCache.TimeToLive = time.Duration(channel.Ttl)
		}
		feedCache.Title = channel.Title
		for _, item := range channel.Item {
			_, err = memcache.Gob.Get(context, item.Guid, nil)
			if err == memcache.ErrCacheMiss {
				err = nil
				date, err = getDate(item.PubDate)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error105: %v\n", err.Error())
					continue
				}
				article := ArticleCache{
					URL:     item.Link,
					Title:   item.Title,
					Summary: item.Description,
					ID:      item.Guid,
					Date:    date.Unix(),
				}
				err = memcache.Gob.Set(context, &memcache.Item{Key: item.Guid, Object: article})
				if err != nil {
					fmt.Fprintf(os.Stderr, "error5: %v\n", err.Error())
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

func getAtom(context appengine.Context, body []byte) (feedCache FeedCache, err error) {
	var feed AtomFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error115: %v\n", err.Error())
	}
	feedCache.Title = feed.Title
	var date time.Time
	for _, item := range feed.Entry {
		_, err = memcache.Gob.Get(context, item.Id, nil)
		if err == memcache.ErrCacheMiss {
			err = nil
			date, err = getDate(item.Updated)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error105: %v\n", err.Error())
				continue
			}
			article := ArticleCache{
				URL:     item.Link[0].Href,
				Title:   item.Title,
				Summary: item.Content.Text,
				ID:      item.Id,
				Date:    date.Unix(),
			}
			err = memcache.Gob.Set(context, &memcache.Item{Key: item.Id, Object: article})
			if err != nil {
				fmt.Fprintf(os.Stderr, "error7: %v\n", err.Error())
				continue
			}
			feedCache.Articles = append(feedCache.Articles, article)
		} else if err != nil {
			break
		}
	}
	return
}

func getSubscription(context appengine.Context, format FeedFormat, body []byte) (feed FeedCache, err error) {
	switch format {
	case RSS:
		return getRSS(context, body)
	case ATOM:
		return getAtom(context, body)
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
	return getSubscription(context, format, body)
}

func getSuggestedFeeds(context appengine.Context, userdata UserData) (suggestedFeeds []string, err error) {
	for _, defaultURL := range defaultFeeds {
		if !ContainsString(suggestedFeeds, defaultURL) && !ContainsString(userdata.Feeds, defaultURL) {
			suggestedFeeds = append(suggestedFeeds, defaultURL)
		}
	}
	query := datastore.NewQuery("Feed")
	var feed Feed
	for iterator := query.Run(context); ; {
		_, err = iterator.Next(&feed)
		if err == datastore.Done {
			err = nil
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "error11: %v\n", err.Error())
			continue
		}
		if !ContainsString(suggestedFeeds, feed.URL) && !ContainsString(userdata.Feeds, feed.URL) {
			suggestedFeeds = append(suggestedFeeds, feed.URL)
		}
	}
	return
}

func feedGET(context appengine.Context, user *user.User, request *http.Request) (data Data, err error) {
	url := request.FormValue("url")
	if url != "" {
		if request.FormValue("unsubscribe") == "1" {
			if user.String() == "default" {
				return
			}
			err = unsubscribe(context, user, url)
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
	feedData.Feeds = userdata.Feeds
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
				fmt.Fprintf(os.Stderr, "error17: %v\n", err.Error())
				continue
			}
		} else {
			err = unsubscribe(context, user, feed.URL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error18: %v\n", err.Error())
				continue
			}
		}
	}
	return
}
