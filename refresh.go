package main

import (
	"crypto/md5"
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
	"appengine/delay"
	"appengine/memcache"
	"appengine/urlfetch"
	"appengine/user"
)

func refreshSubscription(context appengine.Context, feed Feed, feedkey *datastore.Key) (err error) {
	now := time.Now()
	var subscription SubscriptionCache
	var item *memcache.Item
	item, err = memcache.Gob.Get(context, feed.URL, &subscription)
	if err == memcache.ErrCacheMiss {
		err = nil
		subscription.URL = feed.URL
		item = &memcache.Item{Key: feed.URL}
	} else if err != nil {
		return
	}
	if now.Unix() > subscription.Update {
		client := urlfetch.Client(context)
		var response *http.Response
		response, err = client.Get(feed.URL)
		if err != nil {
			return
		}
		defer response.Body.Close()
		var body []byte
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return
		}
		if len(body) != subscription.Length {
			subscription.Length = len(body)
			hash := md5.New()
			var hashLength int
			hashLength, err = hash.Write(body)
			if hashLength != len(body) {
				return errors.New("couldn't make md5 hash")
			}
			if err != nil {
				return
			}
			sum := fmt.Sprintf("%x", hash.Sum(nil))
			if sum != subscription.MD5 {
				if subscription.Format == UNKNOWN {
					subscription.Format = getFeedType(response, body)
				}
				if subscription.Format == OTHER {
					err = memcache.Delete(context, feed.URL)
					if err != nil {
						return
					}
					return errors.New(fmt.Sprintf("%v is not a feed, deleted", feed.URL))
				}
				var articlesCache []ArticleCache
				articlesCache, err = getSubscription(context, subscription.Format, body)
				for _, article := range articlesCache {
					if !ContainsString(feed.Articles, article.ID) {
						feed.Articles = append(feed.Articles, article.ID)
						var articlePrefs = []Pref{
							{
								Field: "feed",
								Value: article.URL,
								Score: 1,
							},
						}
						err = addArticle(context, feed, article.ID, articlePrefs)
						if err != nil {
							fmt.Fprintf(os.Stderr, "error16: %v\n", err.Error())
							continue
						}
					}
				}
				_, err = datastore.Put(context, feedkey, &feed)
				subscription.MD5 = sum
			}
		}
		subscription.Update = time.Now().Add(time.Hour).Unix()
		item.Object = subscription
		err = memcache.Gob.Set(context, item)
	}
	return
}

func refreshSubscriptionURL(context appengine.Context, url string) (err error) {
	query := datastore.NewQuery("Feed").Filter("URL=", url)
	iterator := query.Run(context)
	var feed Feed
	var feedkey *datastore.Key
	feedkey, err = iterator.Next(&feed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error100\n")
		return
	}
	return refreshSubscription(context, feed, feedkey)
}

func refresh(context appengine.Context, asNeeded bool) (data Data, err error) {
	query := datastore.NewQuery("Feed")
	for iterator := query.Run(context); ; {
		var feed Feed
		var feedkey *datastore.Key
		feedkey, err = iterator.Next(&feed)
		if err == datastore.Done {
			err = nil
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "error9: %v\n", err.Error())
			continue
		}
		if asNeeded {
			err = refreshSubscription(context, feed, feedkey)
		} else {
			_, err = getSubscriptionURL(context, feed.URL)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "error10: %v\n", err.Error())
			continue
		}
	}
	return
}

func refreshGET(context appengine.Context, user *user.User, request *http.Request) (data Data, err error) {
	url := request.FormValue("url")
	if url != "" {
		return nil, refreshSubscriptionURL(context, url)
	}
	force := request.FormValue("force")
	if force == "1" {
		return refresh(context, false)
	}
	return refresh(context, true)
}

var refreshDelay = delay.Func("refresh", refreshSubscriptionURL)