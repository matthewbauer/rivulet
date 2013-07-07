package main

import (
	"encoding/xml"
)

type Opml struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    OpmlHead
	Body    OpmlBody
}

type OpmlHead struct {
	XMLName         xml.Name `xml:"head"`
	Title           string   `xml:"title"`
	DateCreated     string   `xml:"dateCreated"`
	DateModified    string   `xml:"dateModified"`
	OwnerName       string   `xml:"ownerName"`
	OwnerEmail      string   `xml:"ownerEmail"`
	ExpansionState  string   `xml:"expansionState"`
	VertScrollState string   `xml:"vertScrollState"`
	WindowTop       string   `xml:"windowTop"`
	WindowLeft      string   `xml:"windowLeft"`
	WindowBottom    string   `xml:"windowBottom"`
	WindowRight     string   `xml:"windowRight"`
}
type OpmlBody struct {
	XMLName  xml.Name      `xml:"body"`
	Outlines []OpmlOutline `xml:"outline"`
}

type OpmlOutline struct {
	XMLName      xml.Name      `xml:"outline"`
	Text         string        `xml:"text,attr,omitempty"`
	Title        string        `xml:"title,attr,omitempty"`
	Type         string        `xml:"type,attr,omitempty"`
	XmlUrl       string        `xml:"xmlUrl,attr,omitempty"`
	HtmlUrl      string        `xml:"htmlUrl,attr,omitempty"`
	IsComment    string        `xml:"isComment,attr,omitempty"`
	IsBreakpoint string        `xml:"isBreakpoint,attr,omitempty"`
	Outlines     []OpmlOutline `xml:"outline"`
}

func getOPMLFeeds(opmlFile []byte) (feeds []string, err error) {
	var opml Opml
	err = xml.Unmarshal(opmlFile, &opml)
	if err != nil {
		return nil, err
	}
	for _, feed := range opml.Body.Outlines {
		feeds = append(feeds, feed.HtmlUrl)
	}
	return
}
