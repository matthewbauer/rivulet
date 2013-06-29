package main

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"os"
	"strings"
	"text/template"
)

import (
	"appengine"
	"appengine/user"
)

func ContainsFeed(list []Feed, elem string) bool {
	for _, t := range list {
		if t.URL == elem {
			return true
		}
	}
	return false
}

func ContainsString(list []string, elem string) bool {
	for _, t := range list {
		if t == elem {
			return true
		}
	}
	return false
}

var templates *template.Template

func init() {
	http.HandleFunc("/", server)
	templates = template.Must(template.ParseFiles("templates/landing.html", "templates/articles.html", "templates/feeds.html", "templates/user.html"))
}

type Data interface {
	Template() string
	Redirect() string
	Send() bool
}

type Redirect struct {
	URL string
}

func (redirect Redirect) Template() string { return "" }
func (redirect Redirect) Redirect() string { return redirect.URL }
func (redirect Redirect) Send() bool       { return true }

type MethodHandler func(appengine.Context, *user.User, *http.Request) (data Data, err error)

var handlers = map[string]map[string]MethodHandler{
	"/_ah/warmup": {
		"GET": warmupGET,
	},
	"/article": {
		"GET":  articleGET,
		"POST": articlePOST,
	},
	"/feed": {
		"GET":  feedGET,
		"POST": feedPOST,
	},
	"/user": {
		"GET":  userGET,
		"POST": userPOST,
	},
	"/refresh": {
		"GET": refreshGET,
	},
	"/app": {
		"GET": appGET,
	},
	"/login": {
		"GET": loginGET,
	},
	"/logout": {
		"GET": logoutGET,
	},
	"/": {
		"GET": rootGET,
	},
}

type OUTPUT int

const (
	UNKNOWNOUTPUT OUTPUT = iota
	JSON
	HTML
)

func mimetypeToOutput(mimetype string) OUTPUT {
	switch mimetype {
	case "application/json", "text/json":
		return JSON
	case "text/html":
		return HTML
	}
	return UNKNOWNOUTPUT
}

func getOutput(request *http.Request) (output OUTPUT) {
	values := []string{request.FormValue("output")}
	accepts := strings.Split(request.Header.Get("Accept"), ",")
	values = append(values, accepts...)
	for _, value := range values {
		mimetype, _, _ := mime.ParseMediaType(value)
		output = mimetypeToOutput(mimetype)
		if output != UNKNOWNOUTPUT {
			return
		}
		mimetype = mime.TypeByExtension(fmt.Sprintf(".%v", value))
		output = mimetypeToOutput(mimetype)
		if output != UNKNOWNOUTPUT {
			return
		}
	}
	return
}

func server(writer http.ResponseWriter, request *http.Request) {
	var err error
	context := appengine.NewContext(request)
	output := getOutput(request)
	if handlers[request.URL.Path] == nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "%v: not found", http.StatusNotFound)
		return
	}
	if request.Method == "OPTIONS" {
		var methods []string
		for method := range handlers[request.URL.Path] {
			methods = append(methods, method)
		}
		writer.Header().Set("Allow", strings.Join(methods, ", "))
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	u := user.Current(context)
	var data Data
	if handlers[request.URL.Path][request.Method] == nil {
		if handlers[request.URL.Path]["*"] != nil {
			data, err = handlers[request.URL.Path]["*"](context, u, request)
		} else {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(writer, "%v: method not allowed", http.StatusMethodNotAllowed)
			return
		}
	} else {
		data, err = handlers[request.URL.Path][request.Method](context, u, request)
	}
	if err != nil {
		printError(context, err, "handler")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if data == nil || !data.Send() {
		writer.WriteHeader(http.StatusOK)
		return
	}
	redirect := data.Redirect()
	if redirect != "" && output != JSON {
		writer.Header().Set("Location", redirect)
		writer.WriteHeader(http.StatusFound)
		return
	}
	err = writeOutput(writer, data, output)
}

func writeOutput(writer http.ResponseWriter, data Data, output OUTPUT) (err error) {
	switch output {
	case JSON:
		var bytes []byte
		bytes, err = json.Marshal(data)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(bytes)
	default:
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		writer.WriteHeader(http.StatusOK)
		err = templates.ExecuteTemplate(writer, data.Template(), data)
	}
	return
}

func logoutGET(context appengine.Context, u *user.User, request *http.Request) (data Data, err error) {
	if u != nil {
		var url string
		request.URL.Path = "/"
		url, err = user.LogoutURL(context, request.URL.String())
		if err != nil {
			return
		}
		var redirect Redirect
		redirect.URL = url
		return redirect, nil
	}
	var redirect Redirect
	redirect.URL = "/"
	return redirect, nil
}

func loginGET(context appengine.Context, u *user.User, request *http.Request) (data Data, err error) {
	if u == nil {
		var url string
		request.URL.Path = "/app"
		url, err = user.LoginURL(context, request.URL.String())
		if err != nil {
			return
		}
		var redirect Redirect
		redirect.URL = url
		return redirect, nil
	}
	var redirect Redirect
	redirect.URL = "/app"
	return redirect, nil
}

type LandingData struct {}
func (LandingData) Template() string { return "landing.html" }
func (LandingData) Redirect() string { return "" }
func (LandingData) Send() bool       { return true }

func rootGET(context appengine.Context, user *user.User, request *http.Request) (data Data, err error) {
	if user != nil {
		var redirect Redirect
		redirect.URL = "/app"
		return redirect, nil
	}
	var landingData LandingData
	return landingData, nil
}

func appGET(context appengine.Context, user *user.User, request *http.Request) (data Data, err error) {
	return article(context, user, request, 0)
}

func warmupGET(context appengine.Context, user *user.User, request *http.Request) (data Data, err error) {
	return
}

func printError(context appengine.Context, err error, info string) {
	fmt.Fprintf(os.Stderr, "(%v) %v\n", info, err.Error())
	context.Errorf("(%v) %v\n", info, err.Error())
}

func printInfo(context appengine.Context, info string) {
	fmt.Fprintf(os.Stderr, "%v\n", info)
	context.Infof("%v\n", info)
}
