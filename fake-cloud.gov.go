package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

type LoginPageContext struct {
	QueryArgs map[string]string
}

type ServerConfig struct {
	CallbackUrl *url.URL
}

func RenderLoginPage(w http.ResponseWriter, context *LoginPageContext) {
	data, err := Asset("data/login.html")
	if err != nil {
		panic("Couldn't find login.html!")
	}
	s := string(data)
	t, _ := template.New("login.html").Parse(s)
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, context)
}

func RedirectToCallback(w http.ResponseWriter, u url.URL, code string, state string) {
	q := u.Query()
	q.Set("code", code)
	q.Set("state", state)
	u.RawQuery = q.Encode()
	w.Header().Set("Location", u.String())
	w.WriteHeader(302)
}

func NewHandler(config *ServerConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		BaseHandler(config, w, r)
	}
}

func BaseHandler(config *ServerConfig, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/oauth/authorize" {
		rq := r.URL.Query()
		email := rq.Get("email")
		if len(email) == 0 {
			queryArgs := make(map[string]string)
			for k, v := range r.URL.Query() {
				queryArgs[k] = v[0]
			}
			RenderLoginPage(w, &LoginPageContext{QueryArgs: queryArgs})
		} else {
			RedirectToCallback(w, *config.CallbackUrl, email, rq.Get("state"))
		}
	} else if r.URL.Path == "/oauth/token" {
		// TODO: Finish implementing this based on
		// https://github.com/18F/calc/blob/develop/fake_uaa_provider/views.py
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "TODO: Implement this!")
	} else if r.URL.Path == "/fake-cloud.gov.svg" {
		data, err := Asset("data/fake-cloud.gov.svg")
		if err != nil {
			panic("Couldn't find fake-cloud.gov.svg!")
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		fmt.Fprintf(w, "%s", data)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(404)
		fmt.Fprintf(w, "Not Found")
	}
}

func Urlify(uStr string) *url.URL {
	u, err := url.Parse(uStr)

	if err != nil {
		panic(fmt.Sprintf("'%s' is not a valid URL!", uStr))
	}

	return u
}

func main() {
	print("Listening on port 8080.\n")

	handler := NewHandler(&ServerConfig{
		// TODO: Get this from the environment and/or command line.
		CallbackUrl: Urlify("http://localhost:8000/callback"),
	})

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
