
// +build !appengine

package main

import (
	"net/http"
	"errors"
	"net/url"
)

import (
  "github.com/bradfitz/gomemcache/memcache"
)

type Context int64

func NewContext(r *http.Request) Context {
	return 0
}

// user

type User struct {
  Name string
}

func (u *User) String() string {
	return u.Name
}

func Current(context Context) *User {
	return &User{"default"}
}

func LoginURL(context Context, dest string) (string, error) {
	return "", nil
}

func LogoutURL(context Context, dest string) (string, error) {
	return "", nil
}

// datastore

var Done = errors.New("datastore: query has no more results")
type Key int64

func GetFirst(context Context, kind string, filter string, value string, dst interface{}) (*Key, error){
}

func NewIncompleteKey(context Context, kind string, parent *Key) *Key {
}

func GetAllKeys(context Context, kind string) []*Key {
}

func Get(context Context, key *Key, dst interface{}) error {
}

func Put(context Context, key *Key, dst interface{}) (*Key, error) {
}

// taskqueue

func Run(context Context, path string, params url.Values) { // hacky, we should make a delay
	if path == "/refresh" {
		if parmas.Get("url") != "" {
			refreshSubscriptionURL(context, params.Get("url"))
		} else {
			refresh(context, params.Get("force") != "")
		}
	}
}

// memcache

var ErrCacheMiss = memcache.ErrCacheMiss

type Item memcache.Item

func memcacheGet(context Context, key string) (*Item, error) {
	item, err := memcache.Get(key)
	return (*Item)(item), err
}

func GobGet(context Context, key string, v interface{}) (*Item, error) {
	item, err := memcache.Gob.Get(key, v)
	return (*Item)(item), err
}

func Set(context Context, item *Item) error {
	return memcache.Set((*memcache.Item)(item))
}

func GobSet(context Context, item *Item) error {
	return memcache.Gob.Set((*memcache.Item)(item))
}

func Delete(context Context, key string) error {
	return memcache.Delete(key)
}
