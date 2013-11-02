package main

import (
	"encoding/xml"
	"time"
	"fmt"
)

import (
	"appengine"
	"appengine/memcache"
)

// Even though we hate it, there are much more RSS feeds out there than Atom feeds
// We'll use http://www.rssboard.org/rss-specification

type RSSStruct struct {
	Version string       `xml:"version,attr"`
	Channel []RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	// required
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Item        []RSSItem `xml:"item"`

	// optional
	Language       string       `xml:"language"`
	Copyright      string       `xml:"copyright"`
	ManagingEditor string       `xml:"managingeditor"`
	WebMaster      string       `xml:"webmaster"`
	PubDate        string       `xml:"pubDate"`
	LastBuildDate  string       `xml:"lastBuildDate"`
	Category       string       `xml:"category"`
	Generator      string       `xml:"generator"`
	Docs           string       `xml:"docs"`
	Cloud          string       `xml:"cloud"`
	Ttl            int          `xml:"ttl"`
	Image          RSSImage     `xml:"image"`
	Rating         string       `xml:"rating"`
	TextInput      RSSTextInput `xml:"textInput"`
	SkipHours      string       `xml:"skipHours"`
	SkipDays       string       `xml:"skipDays"`
}

type RSSSource struct {
	Url  string `xml:"url,attr"`
	Text string `xml:",chardata"`
}

type RSSEnclosure struct {
	Url    string `xml:"url,attr"`
	Length string `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

type RSSItem struct {
	Title       string       `xml:"title"`
	Link        string       `xml:"link"`
	Description string       `xml:"description"`
	Author      string       `xml:"author"`
	Category    string       `xml:"category"`
	Comments    string       `xml:"comments"`
	Enclosure   RSSEnclosure `xml:"enclosure"`
	Guid        string       `xml:"guid"`
	PubDate     string       `xml:"pubDate"`
	DCDate      string       `xml:"date"`
	Source      RSSSource    `xml:"source"`
	Content     string       `xml:"encoded"`
}

type RSSTextInput struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Name        string `xml:"name"`
	Link        string `xml:"link"`
}

type RSSImage struct {
	Url   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`

	Width       int    `xml:"width"`
	Height      int    `xml:"height"`
	Description string `xml:"description"`
}

type RSSCloud struct {
	Domain            string `xml:"domain,attr"`
	Port              string `xml:"port,attr"`
	Path              string `xml:"path,attr"`
	RegisterProcedure string `xml:"registerprocedure,attr"`
	Protocol          string `xml:"protocol,attr"`
}

func getRSS(context appengine.Context, body []byte, url string) (feedCache FeedCache, err error) {
	feedCache.URL = url
	var rss RSSStruct
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		//printError(context, err, url)
		err = nil
	}
	var date time.Time
	for _, channel := range rss.Channel {
		if channel.Ttl > 0 {
			feedCache.TimeToLive = time.Duration(channel.Ttl)
		}
		feedCache.Title = channel.Title
		for _, item := range channel.Item {
			if item.Guid == "" {
				break
			}
			_, err = memcache.Gob.Get(context, item.Guid, nil)
			if err == ErrCacheMiss {
				err = nil
				if item.DCDate != "" {
					item.PubDate = item.DCDate
				}
				date, err = getDate(item.PubDate)
				if err != nil {
					printError(context, fmt.Errorf("rss feed %v has dates that look like %v", channel.Link, item.PubDate), url)
					err = nil
					continue
				}
				var content string
				if item.Content != "" {
					content = item.Content
				} else {
					content = item.Description
				}
				article := ArticleCache{
					URL:     item.Link,
					Title:   item.Title,
					Content: content,
					ID:      item.Guid,
					Date:    date.Unix(),
					FeedName:feedCache.Title,
					FeedURL: feedCache.URL,
				}
				err = memcache.Gob.Set(context, &memcache.Item{Key: item.Guid, Object: article})
				if err != nil {
					printError(context, err, url)
					err = nil
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
