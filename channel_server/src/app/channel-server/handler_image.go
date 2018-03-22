package main

import (
	"net/http"
	"strconv"
	"time"

	"channelling"

	"github.com/gorilla/mux"
)

func makeImageHandler(buddyImages channelling.ImageCache, expires time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		image := buddyImages.Get(vars["imageid"])
		if image == nil {
			http.Error(w, "Unknown image", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", image.MimeType())
		w.Header().Set("ETag", image.LastChangeID())
		age := time.Now().Sub(image.LastChange())

		if age >= time.Second {
			w.Header().Set("Age", strconv.Itoa(int(age.Seconds())))
		}

		if expires >= time.Second {
			w.Header().Set("Expires", time.Now().Add(expires).Format(time.RFC1123))
			w.Header().Set("Cache-Control", "public, no-transform, max-age="+strconv.Itoa(int(expires.Seconds())))
		}

		http.ServeContent(w, r, "", image.LastChange(), image.Reader())
	}
}
