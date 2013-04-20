package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"appengine/user"
)

func articlePOST(context appengine.Context, user *user.User, request *http.Request) (data Data, err error) {
	var body []byte
	body, err = ioutil.ReadAll(request.Body)
	if err != nil {
		return
	}
	var articleList ArticleList
	err = json.Unmarshal(body, &articleList)
	if err != nil {
		return
	}
	var userkey *datastore.Key
	var userdata UserData
	userkey, userdata, err = getUserData(context, user.String())
	if err != nil {
		return
	}
	for _, article := range articleList.Articles {
		found := false
		var n int
		for i, a := range userdata.Articles {
			if a.ID == article.ID {
				found = true
				n = i
				break
			}
		}
		if found {
			if article.Interested {
				err = selected(context, userkey, userdata, article)
				if err != nil {
					return
				}
			}
			if article.Read {
				userdata.Articles = userdata.Articles[:n+copy(userdata.Articles[n:], userdata.Articles[n+1:])]
			} else {
				userdata.Articles[n] = article
			}
		} else {
			userdata.Articles = append(userdata.Articles, article)
		}
	}
	_, err = putUserData(context, userkey, userdata)
	return
}

func article(context appengine.Context, user *user.User, request *http.Request, limit int) (data Data, err error) {
	var userdata UserData
	_, userdata, err = getUserData(context, user.String())
	if err != nil {
		return
	}
	if limit > len(userdata.Articles) {
		var redirect Redirect
		redirect.URL = "/feed"
		return redirect, nil
	}
	var articleData ArticleData
	if user.ID != "default" {
		articleData.User = user.String()
	}
	for _, article := range userdata.Articles[0:limit] {
		var articleCache ArticleCache
		_, err = memcache.Gob.Get(context, article.ID, &articleCache)
		if err == memcache.ErrCacheMiss || articleCache.ID != article.ID {
			err = nil
			var feedArticles []ArticleCache
			feedArticles, err = getSubscriptionURL(context, article.Feed)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error14: %v\n", err.Error())
				continue
			}
			for _, articleCache = range feedArticles {
				if articleCache.ID == article.ID {
					break
				}
			}
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "error15: %v\n", err.Error())
			continue
		}
		articleData.Articles = append(articleData.Articles, articleCache)
	}
	return articleData, nil
}

func getArticleById(articles []Article, id string) (article Article) {
	for _, article = range articles {
		if article.ID == id {
			return article
		}
	}
	return
}

func articleGET(context appengine.Context, user *user.User, request *http.Request) (data Data, err error) {
	id := request.FormValue("id")
	if id != "" && user.String() != "default" {
		var userkey *datastore.Key
		var userdata UserData
		userkey, userdata, err = getUserData(context, user.String())
		if err != nil {
			return
		}
		article := getArticleById(userdata.Articles, id)
		err = selected(context, userkey, userdata, article)
		if err != nil {
			return
		}
		var articleCache ArticleCache
		_, err = memcache.Gob.Get(context, id, &articleCache)
		if err != nil {
			return
		}
		return
	}
	requestNumber := request.FormValue("number")
	var number int
	if requestNumber == "" {
		number = 1
	} else {
		number, err = strconv.Atoi(requestNumber)
		if err != nil {
			return
		}
	}
	return article(context, user, request, number)
}

func addArticle(context appengine.Context, feed Feed, id string, articlePrefs []Pref) (err error) {
	for _, subscriber := range feed.Subscribers {
		userkey, userdata, err := getUserData(context, subscriber)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error2: %v\n", err.Error())
			continue
		}
		rank := getRank(articlePrefs, userdata.Prefs)
		article := Article{Feed: feed.URL, ID: id, Rank: rank, Read: false}
		userdata.Articles = append(userdata.Articles, Article{})
		copy(userdata.Articles[1:], userdata.Articles[0:])
		userdata.Articles[0] = article
		_, err = putUserData(context, userkey, userdata)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error4: %v\n", err.Error())
			continue
		}
	}
	return
}
