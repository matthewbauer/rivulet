package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
		if err == nil {
			switch feed.XMLName.Local {
			case "channel", "rss":
				return RSS
			case "feed":
				return ATOM
			default:
				return OTHER
			}
		}
		return OTHER
	}
	return OTHER
}

func getRSS(context appengine.Context, body []byte) (articles []ArticleCache, err error) {
	var rss RSSStruct
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return
	}
	for _, channel := range rss.Channel {
		for _, item := range channel.Item {
			_, err = memcache.Gob.Get(context, item.Guid, nil)
			if err == memcache.ErrCacheMiss {
				err = nil
				article := ArticleCache{
					URL:     item.Link,
					Title:   item.Title,
					Summary: item.Description,
					ID:      item.Guid,
				}
				err = memcache.Gob.Set(context, &memcache.Item{Key: item.Guid, Object: article})
				if err != nil {
					fmt.Fprintf(os.Stderr, "error5: %v\n", err.Error())
					continue
				}
				articles = append(articles, article)
			} else if err != nil {
				break
			}
		}
	}
	return
}

func getAtom(context appengine.Context, body []byte) (articles []ArticleCache, err error) {
	var feed AtomFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return
	}
	for _, item := range feed.Entry {
		_, err = memcache.Gob.Get(context, item.Id, nil)
		if err == memcache.ErrCacheMiss {
			err = nil
			article := ArticleCache{
				URL:     item.Link[0].Href,
				Title:   item.Title,
				Summary: item.Content.Text,
				ID:      item.Id,
			}
			err = memcache.Gob.Set(context, &memcache.Item{Key: item.Id, Object: article})
			if err != nil {
				fmt.Fprintf(os.Stderr, "error7: %v\n", err.Error())
				continue
			}
			articles = append(articles, article)
		} else if err != nil {
			break
		}
	}
	return
}

func getSubscription(context appengine.Context, format FeedFormat, body []byte) (articles []ArticleCache, err error) {
	switch format {
	case RSS:
		return getRSS(context, body)
	case ATOM:
		return getAtom(context, body)
	case OTHER:
		return nil, errors.New(fmt.Sprintf("not a feed")) // later we can delete the feed
	}
	return nil, errors.New(fmt.Sprintf("could not determine format")) // later we can delete the feed
}

func getSubscriptionURL(context appengine.Context, url string) (articles []ArticleCache, err error) {
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

var defaultFeeds = []string{
	"http://scripting.com/rss.xml",
	"http://daringfireball.net/index.xml",
	"http://www.marco.org/rss",
	"http://feeds.feedburner.com/marginalrevolution/feed",
	"http://feeds.feedburner.com/blogspot/MKuf",
	"http://feeds.kottke.org/main",
	"http://feeds.feedburner.com/thebrowser/xrdJ",
	"http://feeds.feedburner.com/538dotcom",
}

func getSuggestedFeeds(context appengine.Context, userdata UserData) (suggestedFeeds []string, err error) {
	for _, defaultURL := range defaultFeeds {
		if !ContainsString(suggestedFeeds, defaultURL) && !ContainsString(userdata.Feeds, defaultURL) {
			suggestedFeeds = append(suggestedFeeds, defaultURL)
		}
	}
	query := datastore.NewQuery("Feed")
	for iterator := query.Run(context); ; {
		var feed Feed
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
			err = subscribe(context, user, url)
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
	_, userdata, err = getUserData(context, user.String())
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
	url := request.FormValue("url")
	if url != "" {
		err = subscribe(context, user, url)
		if err != nil {
			return
		}
		var redirect Redirect
		redirect.URL = "/feed"
		return redirect, nil
	}
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
			err = subscribe(context, user, feed.URL)
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
