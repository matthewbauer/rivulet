package main

import "encoding/xml"

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
	AtomLink    []AtomLink            `xml:"atom:link"`
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
