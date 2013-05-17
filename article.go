package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
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

func (ArticleData) Template() string { return "articles.html" }
func (ArticleData) Redirect() string { return "" }
func (ArticleData) Send() bool       { return true }

type ArticleList struct {
	Articles []Article
}

func (ArticleList) Template() string { return "articles.html" }
func (ArticleList) Redirect() string { return "" }
func (ArticleList) Send() bool       { return true }

func (al ArticleList) Len() int           { return len(al.Articles) }
func (al ArticleList) Swap(i, j int)      { al.Articles[i], al.Articles[j] = al.Articles[j], al.Articles[i] }
func (al ArticleList) Less(i, j int) bool { return al.Articles[i].Rank > al.Articles[j].Rank }

type Article struct { // child of User
	Rank       int64
	Feed       string
	ID         string
	Read       bool
	Interested bool
}

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
	for _, article := range articleList.Articles {
		var n int
		var a Article
		found := false
		for n, a = range userdata.Articles {
			if a.ID == article.ID {
				found = true
				break
			}
		}
		if found {
			//			if article.Interested {
			//				userdata, err = selected(context, userdata, article)
			//				if err != nil {
			//					printError(context, err)
			//				}
			//			}
			//			if read || article.Read {
			//				//userdata.Articles = append(userdata.Articles[:n], userdata.Articles[n+1:]...)
			//				userdata.Articles[n], userdata.Articles = userdata.Articles[len(userdata.Articles)-1], userdata.Articles[:len(userdata.Articles)-1]
			//			} else {
			userdata.Articles[n] = article
			//			}
		} else {
			userdata.Articles = append(userdata.Articles, article)
		}
	}
	_, err = putUserData(context, userkey, userdata)
	return
}

func article(context appengine.Context, user *user.User, request *http.Request, limit int) (data Data, err error) {
	var userdata UserData
	var userkey *datastore.Key
	userkey, userdata, err = mustGetUserData(context, user.String())
	if err != nil {
		return
	}
	if limit > len(userdata.Articles) {
		refreshDelay.Call(context, "")
		var redirect Redirect
		redirect.URL = "/feed"
		return redirect, nil
	}
	var articleData ArticleData
	if user.ID != "default" {
		articleData.User = user.String()
	}
	var articleCache ArticleCache
	var feedCache FeedCache
	for _, article := range userdata.Articles[0:limit] {
		_, err = memcache.Gob.Get(context, article.ID, &articleCache)
		if err == memcache.ErrCacheMiss { // || articleCache.ID != article.ID
			feedCache, err = getSubscriptionURL(context, article.Feed)
			if err != nil {
				printError(context, err)
				continue
			}
			for _, articleCache = range feedCache.Articles {
				if articleCache.ID == article.ID {
					break
				}
			}
		} else if err != nil {
			printError(context, err)
			continue
		}
		articleData.Articles = append(articleData.Articles, articleCache)
	}
	if len(articleData.Articles) < limit {
		var redirect Redirect
		redirect.URL = "/feed"
		return redirect, nil
	}
	userdata.Articles = userdata.Articles[limit:]
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
	if id != "" && user.String() != "default" {
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
	for _, subscriber := range feed.Subscribers {
		userkey, userdata, err := getUserData(context, subscriber)
		if err != nil {
			printError(context, err)
			continue
		}
		article.Rank = articleCache.Date - time.Now().Unix() //getRank(articlePrefs, userdata.Prefs)
		userdata.Articles = append(userdata.Articles, article)
		_, err = putUserData(context, userkey, userdata)
		if err != nil {
			printError(context, err)
			continue
		}
	}
	return
}
