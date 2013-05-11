package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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

var refreshSubscriptionURLDelay = delay.Func("refresh", refreshSubscriptionURL)
var refreshDelay = delay.Func("refresh", func(context appengine.Context, x string) { refresh(context, true) })

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
		duration := time.Hour
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
				printError(context, errors.New("Refreshing..."))
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
				var feedCache FeedCache
				feedCache, err = getSubscription(context, subscription.Format, body)
				if err != nil {
					return
				}
				if feedCache.TimeToLive > 0 {
					duration = time.Duration(feedCache.TimeToLive)
				}
				for _, article := range feedCache.Articles {
					if !ContainsString(feed.Articles, article.ID) {
						feed.Articles = append(feed.Articles, article.ID)
						err = addArticle(context, feed, article)
						if err != nil {
							printError(context, err)
							continue
						}
					}
				}
				_, err = datastore.Put(context, feedkey, &feed)
				subscription.MD5 = sum
			}
		}
		subscription.Update = time.Now().Add(duration).Unix()
		item.Object = subscription
		err = memcache.Gob.Set(context, item)
		if err != nil {
			return
		}
	}
	return
}

func refreshSubscriptionURL(context appengine.Context, url string) (err error) {
	query := datastore.NewQuery("Feed").Filter("URL=", url)
	iterator := query.Run(context)
	var feed Feed
	var feedkey *datastore.Key
	feedkey, err = iterator.Next(&feed)
	if err == datastore.Done {
		feed.URL = url
		feedkey = datastore.NewIncompleteKey(context, "Feed", nil)
	} else if err != nil {
		return
	}
	return refreshSubscription(context, feed, feedkey)
}

func refresh(context appengine.Context, asNeeded bool) (data Data, err error) {
	query := datastore.NewQuery("Feed")
	var feed Feed
	var feedkey *datastore.Key
	for iterator := query.Run(context); ; {
		feedkey, err = iterator.Next(&feed)
		if err == datastore.Done {
			err = nil
			break
		} else if err != nil {
			printError(context, err)
			continue
		}
		if asNeeded {
			err = refreshSubscription(context, feed, feedkey)
		} else {
			_, err = getSubscriptionURL(context, feed.URL)
		}
		if err != nil {
			printError(context, err)
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
