package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
)

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"appengine/user"
)

type ArticleCache struct {
	Title   string
	Summary string
	URL     string
	ID      string
	Feed    string
	Date    int64
}

type ArticleData struct {
	User     string
	Articles []ArticleCache
}

const MAXARTICLES = 100
const DEFAULTCOUNT = 1

func (ArticleData) Template() string { return "articles.html" }
func (ArticleData) Redirect() string { return "" }
func (ArticleData) Send() bool       { return true }

type Article struct {
	Rank       int64
	Feed       string
	ID         string
	Read       bool
	Interested bool
}

type ArticleList struct {
	Articles []Article
}

func (ArticleList) Template() string { return "articles.html" }
func (ArticleList) Redirect() string { return "" }
func (ArticleList) Send() bool       { return true }

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
	userkey, userdata, err = mustGetUserData(context, user.String())
	if err != nil {
		return
	}
	read := false
	if request.FormValue("read") == "1" {
		read = true
	}
	for _, article := range articleList.Articles {
		n := 0
		var a Article
		found := false
		for n, a = range userdata.Articles {
			if a.ID == article.ID {
				found = true
				break
			}
		}
		if found {
			if article.Interested {
				userdata, err = selected(context, userdata, article)
				if err != nil {
					printError(context, err, article.ID)
				}
			}
			if read || article.Read {
				userdata.Articles = append(userdata.Articles[:n], userdata.Articles[n+1:]...)
				userdata.Articles[n], userdata.Articles = userdata.Articles[len(userdata.Articles)-1], userdata.Articles[:len(userdata.Articles)-1]
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

func article(context appengine.Context, user *user.User, request *http.Request, count int) (data Data, err error) {
	var userdata UserData
	var userkey *datastore.Key
	userkey, userdata, err = mustGetUserData(context, user.String())
	if err != nil {
		return
	}
	var articleData ArticleData
	var articleCache ArticleCache
	var feedCache FeedCache
	if count > len(userdata.Articles) {
		refreshDelay.Call(context, "")
		var redirect Redirect
		redirect.URL = "/feed"
		return redirect, nil
	}
	for _, article := range userdata.Articles[0:count] {
		_, err = memcache.Gob.Get(context, article.ID, &articleCache)
		if err == memcache.ErrCacheMiss {
			feedCache, err = getSubscriptionURL(context, article.Feed)
			if err != nil {
				printError(context, err, article.Feed)
				continue
			}
			for _, articleCache = range feedCache.Articles {
				if articleCache.ID == article.ID {
					break
				}
			}
		} else if err != nil {
			printError(context, err, article.ID)
			continue
		}
		if articleCache.URL != "" {
			articleData.Articles = append(articleData.Articles, articleCache)
		}
	}
	//if count > len(articleData.Articles) {
	//	refreshDelay.Call(context, "")
	//	var redirect Redirect
	//	redirect.URL = "/feed"
	//	return redirect, nil
	//}
	userdata.Articles = userdata.Articles[count:]
	_, err = putUserData(context, userkey, userdata)
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
	if id != "" {
		var userkey *datastore.Key
		var userdata UserData
		userkey, userdata, err = mustGetUserData(context, user.String())
		if err != nil {
			return
		}
		article := getArticleById(userdata.Articles, id)
		userdata, err = selected(context, userdata, article)
		if err != nil {
			return
		}
		_, err = putUserData(context, userkey, userdata)
	}
	url := request.FormValue("url")
	if url != "" {
		var redirect Redirect
		redirect.URL = url
		return redirect, nil
	}
	countStr := request.FormValue("count")
	count := 0
	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			return
		}
	} else {
		count = DEFAULTCOUNT
	}

	return article(context, user, request, count)
}

func addArticle(context appengine.Context, feed Feed, articleCache ArticleCache) (err error) {
	if articleCache.ID == "" {
		return
	}
	/*
	*var articlePrefs = []Pref{
	*    {
	*        Field: "feed",
	*        Value: articleCache.URL,
	*        Score: 1,
	*    },
	*}
	 */
	article := Article{Feed: feed.URL, ID: articleCache.ID, Read: false}
	var userkey *datastore.Key
	var userdata UserData
	if feed.Default {
		feed.Subscribers = append(feed.Subscribers, "default")
	}
	for _, subscriber := range feed.Subscribers {
		userkey, userdata, err = getUserData(context, subscriber)
		if err != nil {
			continue
		}
		article.Rank = articleCache.Date //getRank(articlePrefs, userdata.Prefs)
		if len(userdata.Articles) > MAXARTICLES {
			userdata.Articles = userdata.Articles[0 : MAXARTICLES-1]
		}

		n := sort.Search(len(userdata.Articles), func(i int) bool { return userdata.Articles[i].Rank <= article.Rank })
		if n == -1 {
			userdata.Articles = append(userdata.Articles, article)
		} else {
			userdata.Articles = append(userdata.Articles, Article{})
			copy(userdata.Articles[n+1:], userdata.Articles[n:])
			userdata.Articles[n] = article
		}

		_, err = putUserData(context, userkey, userdata)
		if err != nil {
			printError(context, err, subscriber)
			continue
		}
	}
	err = nil
	return
}
