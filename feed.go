package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"encoding/json"
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
	HTMLFeed
	OTHER
)

type FeedInfo struct {
	URL        string
}

type FeedList struct {
	Feeds []string
}

type FeedData struct {
	User           string
	Feeds          []Feed
	SuggestedFeeds []Feed
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
	Title       string
	Default     bool
	Subscribers []string
	Articles    []string
}

type GenericFeed struct {
	XMLName xml.Name
}

type HTMLLink struct {
	Rel string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
	Title string `xml:"title,attr"`
	Href string `xml:"href,attr"`
}

type HTMLHead struct {
	Link []HTMLLink `xml:"link"`
}

type HTMLStruct struct {
	Head HTMLHead `xml:"head"`
}

func getFeedType(response *http.Response, body []byte) FeedFormat {
	switch response.Header.Get("Content-Type") {
	case "application/atom+xml":
		return ATOM
	case "application/rss+xml":
		return RSS
	case "text/html":
		return HTMLFeed
	default:
		var feed GenericFeed
		xml.Unmarshal(body, &feed)
		switch feed.XMLName.Local {
		case "channel", "rss":
			return RSS
		case "feed":
			return ATOM
		case "html":
			return HTMLFeed
		default:
			return OTHER
		}
		return OTHER
	}
	return OTHER
}

const (
	ONIONDATE = "Mon, 2 Jan 2006 15:04:05 -0700"
	MYDATE    = "1/2/2006 3:04:05 PM"
)

func getDate(dateString string) (date time.Time, err error) {
	layouts := []string{time.RFC822, time.RFC822Z, time.RFC3339, time.RFC1123, time.RFC1123Z,
		time.ANSIC, time.UnixDate, time.RubyDate, ONIONDATE} //, MYDATE
	for _, layout := range layouts {
		date, err = time.Parse(layout, dateString)
		if err == nil && date.Year() != 0 {
			return
		}
	}
	return
}

var (
	NotFeed = errors.New(fmt.Sprintf("not a feed"))
	UnknownFormat = errors.New(fmt.Sprintf("could not determine format"))
)

func getSubscription(context appengine.Context, format FeedFormat, body []byte, url string) (feed FeedCache, err error) {
	// url is only used for debugging purposes
	switch format {
	case RSS:
		return getRSS(context, body, url)
	case ATOM:
		return getAtom(context, body, url)
	case OTHER:
		err = NotFeed // later we can delete the feed
		return
	}
	err = UnknownFormat // later we can delete the feed
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

func feedURLGET(context appengine.Context, u string) (data Data, err error) {
	var feed Feed
	_, err = GetFirst(context, "Feed", "URL", u, &feed)

	articleData := ArticleData{ Source: fmt.Sprintf("/feed?url=%v", url.QueryEscape(u) )}
	for _, article := range feed.Articles {
		var articleCache ArticleCache
		_, err = memcache.Gob.Get(context, article, &articleCache)
		if err != nil {
			printError(context, err, feed.URL)
			err = nil
			continue
		}
		if articleCache.URL != "" {
			articleData.Articles = append(articleData.Articles, articleCache)
		}
	}

	return articleData, err
}

func feedGET(context appengine.Context, user *user.User, request *http.Request) (data Data, err error) {
	u := request.FormValue("url")
	if u != "" {
		return feedURLGET(context, u)
	}

	feedData := FeedData{ User: user.String() }

	var userdata UserData
	_, userdata, err = mustGetUserData(context, user.String())
	if err != nil {
		return
	}

	for _, feed := range userdata.Feeds {
		var item Feed
		for _, defaultFeed := range builtinFeeds {
			if defaultFeed.URL == feed {
				item = defaultFeed
				break
			}
		}

		if item.URL == "" {
			_, err = GetFirst(context, "Feed", "URL", feed, &item)
			if err != nil {
				printError(context, err, feed)
				err = nil
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

func subscribeFeedList(context appengine.Context, userdata *UserData, feedList FeedList) (err error) {
	for _, feed := range feedList.Feeds {
		err = subscribe(context, userdata, feed, false)
		if err != nil {
			printError(context, err, feed)
			err = nil
			continue
		}
	}
	return
}

func unsubscribeFeedList(context appengine.Context, userdata *UserData, feedList FeedList) (err error) {
	for _, feed := range feedList.Feeds {
		err = unsubscribe(context, userdata, feed)
		if err != nil {
			printError(context, err, feed)
			err = nil
			continue
		}
	}
	return
}

func unsubscribeAll(context appengine.Context, userdata *UserData) (err error) {
	var feedList FeedList
	for _, feed := range userdata.Feeds {
		feedList.Feeds = append(feedList.Feeds, feed)
	}

	err = unsubscribeFeedList(context, userdata, feedList)
	return
}

func feedJSONPOST(context appengine.Context, userdata *UserData, request *http.Request) (err error) {
	// todo: add json subscribe
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

	err = subscribeFeedList(context, userdata, feedList)
	return
}

func feedPOST(context appengine.Context, user *user.User, request *http.Request) (data Data, err error) {
	var userdata UserData
	var userkey *datastore.Key
	userkey, userdata, err = mustGetUserData(context, user.String())
	if err != nil {
		return
	}

	if request.FormValue("clear") != "" {
		err = unsubscribeAll(context, &userdata)
	}

	if request.FormValue("input") == "json" {
		err = feedJSONPOST(context, &userdata, request)
	} else if request.FormValue("input") == "opml" {
		err = feedOPMLPOST(context, &userdata, request)
	} else if request.FormValue("input") == "form" || request.FormValue("url") != "" {
		err = subscribe(context, &userdata, request.FormValue("url"), false)
	}
	if err != nil {
		return
	}

	_, err = putUserData(context, userkey, userdata)
	if err != nil {
		return
	}

	redirect := Redirect{URL: "/app"}
	return redirect, nil
}

func feedDELETE(context appengine.Context, user *user.User, request *http.Request) (data Data, err error) {
	var userdata UserData
	var userkey *datastore.Key
	userkey, userdata, err = mustGetUserData(context, user.String())
	if err != nil {
		return
	}

	if request.FormValue("url") != "" {
		err = unsubscribe(context, &userdata, request.FormValue("url"))
	} else {
		err = unsubscribeAll(context, &userdata)
	}
	if err != nil {
		return
	}

	_, err = putUserData(context, userkey, userdata)
	if err != nil {
		return
	}

	redirect := Redirect{URL: "/app"}
	return redirect, nil
}
