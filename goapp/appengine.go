
// +build appengine

package main

import (
	"net/http"
	"net/url"
	"bytes"
	"fmt"
	"encoding/gob"
	"text/template"
)

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"appengine/user"
	"appengine/taskqueue"
)

func init() {
	http.HandleFunc("/", server)
	templates = template.Must(template.ParseFiles("templates/landing.html", "templates/articles.html", "templates/feeds.html", "templates/user.html"))
}

type Context appengine.Context
func NewContext(r *http.Request) Context {
	return appengine.NewContext(r)
}

// user

type User user.User

func (u *User) String() string {
	return u.String()
}

func Current(context Context) *User {
	return (*User)(user.Current(context))
}

func LoginURL(context Context, dest string) (string, error) {
	return user.LoginURL(context, dest)
}

func LogoutURL(context Context, dest string) (string, error) {
	return user.LogoutURL(context, dest)
}

// datastore

func (userdata *UserData) Load(c <-chan datastore.Property) (err error) {
	for p := range c {
		switch p.Name {
		case "String":
			userdata.String = p.Value.(string)
		case "TotalRead":
			userdata.TotalRead = p.Value.(int64)
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
	c <- datastore.Property{
		Name:  "TotalRead",
		Value: userdata.TotalRead,
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

var Done = datastore.Done
type Key datastore.Key

func GetFirst(context Context, kind string, filter string, value string, dst interface{}) (*Key, error){
	query := datastore.NewQuery(kind).Filter(fmt.Sprintf("%v=", filter), value)
	iterator := query.Run(appengine.Context(context))
	feedkey, err := iterator.Next(dst)
	return (*Key)(feedkey), err
}

func NewIncompleteKey(context Context, kind string, parent *Key) *Key {
	return (*Key)(datastore.NewIncompleteKey(appengine.Context(context), kind, (*datastore.Key)(parent)))
}

func GetAllKeys(context Context, kind string) (retkeys []*Key, err error) {
	query := datastore.NewQuery(kind)
	var keys []*datastore.Key
	keys, err = query.KeysOnly().GetAll(context, nil)
	for _, key := range keys {
		retkeys = append(retkeys, (*Key)(key))
	}
	return
}

func Get(context Context, key *Key, dst interface{}) error {
	return datastore.Get(appengine.Context(context), (*datastore.Key)(key), dst)
}

func Put(context Context, key *Key, dst interface{}) (*Key, error) {
	newkey, err := datastore.Put(appengine.Context(context), (*datastore.Key)(key), dst)
	return (*Key)(newkey), err
}

// delay

func Run(context Context, path string, params url.Values) {
	t := taskqueue.NewPOSTTask(path, params)
	taskqueue.Add(appengine.Context(context), t, "")
}

// memcache

var ErrCacheMiss = memcache.ErrCacheMiss

type Item memcache.Item

func memcacheGet(context Context, key string) (*Item, error) {
	item, err := memcache.Get(appengine.Context(context), key)
	return (*Item)(item), err
}

func GobGet(context Context, key string, v interface{}) (*Item, error) {
	item, err := memcache.Gob.Get(appengine.Context(context), key, v)
	return (*Item)(item), err
}

func Set(context Context, item *Item) error {
	return memcache.Set(appengine.Context(context), (*memcache.Item)(item))
}

func GobSet(context Context, item *Item) error {
	return memcache.Gob.Set(appengine.Context(context), (*memcache.Item)(item))
}

func Delete(c appengine.Context, key string) error {
	return memcache.Delete(appengine.Context(c), key)
}
