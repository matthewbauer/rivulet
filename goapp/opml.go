package main

import (
	"encoding/xml"
)

type Opml1 struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    OpmlHead
	Body    OpmlBody
}

type Opml2 struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    OpmlHead2
	Body    OpmlBody2
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

type OpmlHead2 struct {
	XMLName         xml.Name `xml:"head"`
	Title           string   `xml:"title"`
	DateCreated     string   `xml:"dateCreated"`
	DateModified    string   `xml:"dateModified"`
	OwnerName       string   `xml:"ownerName"`
	OwnerEmail      string   `xml:"ownerEmail"`
	OwnerId         string   `xml:"ownerId"`
	Docs            string   `xml:"docs"`
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

type OpmlBody2 struct {
	XMLName  xml.Name      `xml:"body"`
	Outlines []OpmlOutline `xml:"outline"`
}

type OpmlOutline2 struct {
	XMLName      xml.Name      `xml:"outline"`
	Text         string        `xml:"text,attr,omitempty"`
	Title        string        `xml:"title,attr,omitempty"`
	Type         string        `xml:"type,attr,omitempty"`
	XmlUrl       string        `xml:"xmlUrl,attr,omitempty"`
	HtmlUrl      string        `xml:"htmlUrl,attr,omitempty"`
	IsComment    string        `xml:"isComment,attr,omitempty"`
	IsBreakpoint string        `xml:"isBreakpoint,attr,omitempty"`
	Created      string        `xml:"created,attr,omitempty"`
	Category     string        `xml:"category,attr,omitempty"`
	Description  string        `xml:"description,attr,omitempty"`
	Language     string        `xml:"language,attr,omitempty"`
	Version      string        `xml:"version,attr,omitempty"`
	Url          string        `xml:"url,attr,omitempty"`
	Outlines     []OpmlOutline `xml:"outline"`
}


