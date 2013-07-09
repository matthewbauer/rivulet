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

// Based on rfc4287
type AtomCategory struct {
	Term   string
	Scheme string
	Label  string
}

type AtomContent struct {
	Data string `xml:"data,attr"`
	Type string `xml:"type,attr"`
	Text string `xml:",chardata"`
}

type AtomEntry struct {
	Author      []AtomPersonConstruct `xml:"author"`
	Content     AtomContent           `xml:"content"`
	Category    []AtomCategory        `xml:"category"`
	Contributor []AtomPersonConstruct `xml:"contributor"`
	Id          string                `xml:"id"`
	Link        []AtomLink            `xml:"link"`
	Published   string                `xml:"published"`
	Rights      string                `xml:"rights"`
	Source      AtomSource            `xml:"source"`
	Summary     string                `xml:"summary"`
	Title       string                `xml:"title"`
	Updated     string                `xml:"updated"`
}

type AtomFeed struct {
	XMLName     xml.Name              `xml:"http://www.w3.org/2005/Atom feed"`
	Author      []AtomPersonConstruct `xml:"author"`
	Category    []AtomCategory        `xml:"category"`
	Contributor []AtomPersonConstruct `xml:"contributor"`
	Entry       []AtomEntry           `xml:"entry"`
	Generator   AtomGenerator         `xml:"generator"`
	Icon        string                `xml:"icon"`
	Id          string                `xml:"id"`
	Link        []AtomLink            `xml:"link"`
	Logo        string                `xml:"logo"`
	Rights      string                `xml:"rights"`
	Subtitle    string                `xml:"subtitle"`
	Title       string                `xml:"title"`
	Updated     string                `xml:"updated"`
}

type AtomGenerator struct {
	Uri     string `xml:"uri,attr"`
	Version string `xml:"version,attr"`
	Text    string `xml:",chardata"`
}

type AtomLink struct {
	Href     string `xml:"href,attr"`
	Rel      string `xml:"rel,attr"`
	Type     string `xml:"type,attr"`
	Hreflang string `xml:"hreflang,attr"`
	Title    string `xml:"title,attr"`
	Length   string `xml:"length,attr"`
}

type AtomPersonConstruct struct {
	Email string `xml:"email"`
	Name  string `xml:"name"`
	Uri   string `xml:"uri"`
}

type AtomSource struct {
	Author      []AtomPersonConstruct `xml:"author"`
	Category    []AtomCategory        `xml:"category"`
	Contributor []AtomPersonConstruct `xml:"contributor"`
	Generator   AtomGenerator         `xml:"generator"`
	Icon        string                `xml:"icon"`
	Id          string                `xml:"id"`
	Link        []AtomLink            `xml:"link"`
	Logo        string                `xml:"logo"`
	Rights      string                `xml:"rights"`
	Subtitle    string                `xml:"subtitle"`
	Title       string                `xml:"title"`
	Updated     string                `xml:"updated"`
}

func getAtom(context appengine.Context, body []byte, url string) (feedCache FeedCache, err error) {
	feedCache.URL = url
	var feed AtomFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		//printError(context, err, url)
		err = nil
	}
	feedCache.Title = feed.Title
	var date time.Time
	for _, item := range feed.Entry {
		if item.Id == "" {
			break
		}
		_, err = memcache.Gob.Get(context, item.Id, nil)
		if err == memcache.ErrCacheMiss {
			err = nil
			date, err = getDate(item.Updated)
			if err != nil {
				printError(context, fmt.Errorf("atom feed %v has dates that look like %v", feed.Link[0].Href, item.Updated), url)
				err = nil
				continue
			}
			var url string
			for _, link := range item.Link {
				if link.Href != "" {
					url = link.Href
					if link.Rel == "alternate" {
						break
					}
				}
			}
			if url == "" {
				break
			}
			article := ArticleCache{
				URL:     url,
				Title:   item.Title,
				Content: item.Content.Text,
				ID:      item.Id,
				Date:    date.Unix(),
				FeedName:feedCache.Title,
				FeedURL: feedCache.URL,
			}
			err = memcache.Gob.Set(context, &memcache.Item{Key: item.Id, Object: article})
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
	return
}

