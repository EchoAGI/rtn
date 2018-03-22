package main

import (
	"net/http"

	"channelling"

	"github.com/gorilla/mux"
)

func roomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	handleRoomView(vars["room"], w, r)
}

func handleRoomView(room string, w http.ResponseWriter, r *http.Request) {
	var err error

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Expires", "-1")
	w.Header().Set("Cache-Control", "private, max-age=0")

	csp := false

	if config.ContentSecurityPolicy != "" {
		w.Header().Set("Content-Security-Policy", config.ContentSecurityPolicy)
		csp = true
	}
	if config.ContentSecurityPolicyReportOnly != "" {
		w.Header().Set("Content-Security-Policy-Report-Only", config.ContentSecurityPolicyReportOnly)
		csp = true
	}

	scheme := "http"

	// Detect if the request was made with SSL.
	ssl := r.TLS != nil
	proto, ok := r.Header["X-Forwarded-Proto"]
	if ok {
		ssl = proto[0] == "https"
		scheme = "https"
	}

	// Get languages from request.
	langs := getRequestLanguages(r, []string{})
	if len(langs) == 0 {
		langs = append(langs, "en")
	}

	// Prepare context to deliver to HTML..
	context := &channelling.Context{
		Cfg:        config,
		App:        "main",
		Host:       r.Host,
		Scheme:     scheme,
		Ssl:        ssl,
		Csp:        csp,
		Languages:  langs,
		Room:       room,
		S:          config.S,
		ExtraDHead: templatesExtraDHead,
		ExtraDBody: templatesExtraDBody,
	}

	// Get URL parameters.
	r.ParseForm()

	// Check if incoming request is a crawler which supports AJAX crawling.
	// See https://developers.google.com/webmasters/ajax-crawling/docs/getting-started for details.
	if _, ok := r.Form["_escaped_fragment_"]; ok {
		// Render crawlerPage template..
		err = templates.ExecuteTemplate(w, "crawlerPage", context)
	} else {
		// Render mainPage template.
		err = templates.ExecuteTemplate(w, "mainPage", context)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
