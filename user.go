package main

import (
	"net/http"
)

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func getUserData(context appengine.Context, id string) (key *datastore.Key, userdata UserData, err error) {
	query := datastore.NewQuery("UserData").Filter("String=", id).Limit(1)
	iterator := query.Run(context)
	key, err = iterator.Next(&userdata)
	if err == datastore.Done {
		userdata.String = id
		userdata.Feeds = defaultFeeds
		key, err = putUserData(context, datastore.NewIncompleteKey(context, "UserData", nil), userdata)
	}
	return
}

func putUserData(context appengine.Context, oldkey *datastore.Key, userdata UserData) (newkey *datastore.Key, err error) {
	newkey, err = datastore.Put(context, oldkey, &userdata)
	return
}

func getRank(article []Pref, user []Pref) (score int) {
	for _, userPref := range user {
		for _, articlePref := range article {
			if userPref.Field == articlePref.Field && userPref.Value == articlePref.Value {
				score += articlePref.Score * userPref.Score
			}
		}
	}
	return
}

func unsubscribe(context appengine.Context, user *user.User, url string) (err error) {
	var userdata UserData
	var userkey *datastore.Key
	userkey, userdata, err = getUserData(context, user.String())
	if err != nil {
		return
	}
	for i, feed := range userdata.Feeds {
		if feed == url {
			userdata.Feeds = userdata.Feeds[:i+copy(userdata.Feeds[i:], userdata.Feeds[i+1:])]
			_, err = putUserData(context, userkey, userdata)
			if err != nil {
				return
			}
			break
		}
	}
	query := datastore.NewQuery("Feed").Filter("URL=", url)
	var feed Feed
	var key *datastore.Key
	iterator := query.Run(context)
	key, err = iterator.Next(&feed)
	if err != nil {
		return
	}
	for i, subscriber := range feed.Subscribers {
		if subscriber == userdata.String {
			feed.Subscribers = feed.Subscribers[:i+copy(feed.Subscribers[i:], feed.Subscribers[i+1:])]
			_, err = datastore.Put(context, key, &feed)
			break
		}
	}
	return
}

func subscribe(context appengine.Context, user *user.User, url string) (err error) {
	query := datastore.NewQuery("Feed").Filter("URL=", url)
	iterator := query.Run(context)
	var feed Feed
	var key *datastore.Key
	feedsubscribed := false
	key, err = iterator.Next(&feed)
	if err == datastore.Done {
		feed.URL = url
		feed.Subscribers = []string{user.String()}
		key, err = datastore.Put(context, datastore.NewIncompleteKey(context, "Feed", nil), &feed)
		refreshDelay.Call(context, feed.URL)
		feedsubscribed = true
	}
	var userdata UserData
	var userkey *datastore.Key
	userkey, userdata, err = getUserData(context, user.String())
	if err != nil {
		return
	}
	if !ContainsString(userdata.Feeds, url) {
		userdata.Feeds = append(userdata.Feeds, url)
		_, err = putUserData(context, userkey, userdata)
		if err != nil {
			return
		}
		if !feedsubscribed {
			feed.Subscribers = append(feed.Subscribers, user.String())
			_, err = datastore.Put(context, key, &feed)
		}
	}
	return
}

func selected(context appengine.Context, userkey *datastore.Key, userdata UserData, article Article) (err error) {
	found := false
	for i, value := range userdata.Prefs {
		if value.Field == "field" && value.Value == article.Feed {
			found = true
			value.Score += 1
			userdata.Prefs[i] = value
			break
		}
	}
	if !found {
		userdata.Prefs = append(userdata.Prefs, Pref{
			Field: "feed",
			Value: article.Feed,
			Score: 1,
		})
	}
	_, err = putUserData(context, userkey, userdata)
	return
}

func userGET(context appengine.Context, user *user.User, request *http.Request) (data Data, err error) {
	var userdata UserData
	_, userdata, err = getUserData(context, user.String())
	if err != nil {
		return
	}
	return userdata, nil
}
