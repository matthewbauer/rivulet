
// +build !appengine 

package main

import (
	"net/http"
	"github.com/bradfitz/gomemcache/memcache"
)

// context

type Context struct {
	Memcache *memcache.Client
}

func NewContext(request *http.Request) Context {
	return Context{}
}

// memcache

type Item memcache.Item

func memcacheGet(c Context, key string) (*Item, error) {
	item, err := c.Memcache.Get(key)
	return (*Item)(item), err
}

func memcacheSet(c Context, item *Item) error {
	return c.Memcache.Set((*memcache.Item)(item))
}

func memcacheGobGet(c Context, key string, v interface{}) (*Item, error) {
	item, err := memcache.Gob.Get(appengine.Context(c), key, v)
	return (*Item)(item), err
}

func memcacheGobSet(c Context, item *Item) error {
	return memcache.Gob.Set(appengine.Context(c), (*memcache.Item)(item))
}
