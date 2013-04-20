package main

import (
	"bytes"
	"encoding/gob"
)

import (
	"appengine/datastore"
)

type FeedInfo struct {
	URL        string
	Subscribed bool
}

type FeedList struct {
	Feeds []FeedInfo
}

type ArticleData struct {
	User     string
	Articles []ArticleCache
}

func (ArticleData) Template() string { return "articles.html" }
func (ArticleData) Redirect() string { return "" }
func (ArticleData) Send() bool       { return true }

type FeedData struct {
	User           string
	Feeds          []string
	SuggestedFeeds []string
}

func (FeedData) Template() string { return "feeds.html" }
func (FeedData) Redirect() string { return "" }
func (FeedData) Send() bool       { return true }

type SubscriptionCache struct {
	URL    string
	Format FeedFormat
	Update int64
	MD5    string
	Length int
}

type ArticleList struct {
	Articles []Article
}

func (ArticleList) Template() string { return "articles.html" }
func (ArticleList) Redirect() string { return "" }
func (ArticleList) Send() bool       { return true }

type ArticleCache struct {
	Title   string
	Summary string
	URL     string
	ID      string
}

type Article struct { // child of User
	Rank       int
	Feed       string
	ID         string
	Read       bool
	Interested bool
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

type Pref struct {
	Field string
	Value string
	Score int
}

type UserData struct {
	String   string
	Bytes    []byte    `datastore:",noindex"`
	Articles []Article `datastore:"-"`
	Feeds    []string  `datastore:"-"`
	Prefs    []Pref    `datastore:"-"`
}

func (UserData) Template() string { return "user.html" }
func (UserData) Redirect() string { return "" }
func (UserData) Send() bool       { return true }

func (userdata *UserData) Load(c <-chan datastore.Property) (err error) {
	for p := range c {
		switch p.Name {
		case "String":
			userdata.String = p.Value.(string)
		case "Bytes":
			reader := bytes.NewBuffer(p.Value.([]byte))
			decoder := gob.NewDecoder(reader)
			err = decoder.Decode(userdata)
			if err != nil {
				return
			}
		}
	}
	return
}

func (userdata *UserData) Save(c chan<- datastore.Property) (err error) {
	defer close(c)
	c <- datastore.Property{
		Name:  "String",
		Value: userdata.String,
	}
	writer := bytes.Buffer{}
	encoder := gob.NewEncoder(&writer)
	err = encoder.Encode(userdata)
	if err != nil {
		return
	}
	c <- datastore.Property{Name: "Bytes", Value: writer.Bytes(), NoIndex: true}
	return
}
