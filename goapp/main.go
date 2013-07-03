// +build !appengine

package main

import "net/http"

func main() {
	http.HandleFunc("/", server)
	http.ListenAndServe("localhost:8080", nil)
}
