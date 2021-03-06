
// +build appengine 

package main

import (
	"net/http"
)

import "appengine"
import "appengine/memcache"

// context

type Context appengine.Context

func NewContext(request *http.Request) Context {
	return Context(appengine.NewContext(request))
}

// memcache

var ErrCacheMiss = memcache.ErrCacheMiss

type Item memcache.Item

func memcacheGet(c Context, key string) (*Item, error) {
	item, err := memcache.Get(appengine.Context(c), key)
	return (*Item)(item), err
}

func memcacheSet(c Context, item *Item) error {
	return memcache.Set(appengine.Context(c), (*memcache.Item)(item))
}

func memcacheGobGet(c Context, key string, v interface{}) (*Item, error) {
	item, err := memcache.Gob.Get(appengine.Context(c), key, v)
	return (*Item)(item), err
}

func memcacheGobSet(c Context, item *Item) error {
	return memcache.Gob.Set(appengine.Context(c), (*memcache.Item)(item))
}
