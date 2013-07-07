// +build !appengine

package main

import "net/http"

func main() {
	templates = template.Must(template.ParseFiles("templates/landing.html", "templates/articles.html", "templates/feeds.html", "templates/user.html"))
	http.HandleFunc("/", server)
	http.ListenAndServe("localhost:8080", nil)
}
