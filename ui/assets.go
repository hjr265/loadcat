// Copyright 2015 The Loadcat Authors. All rights reserved.

package ui

import (
	"net/http"
)

func ServeAsset(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path)
}

func init() {
	Router.NewRoute().
		Methods("GET").
		PathPrefix("/assets").
		Handler(http.StripPrefix("/assets", http.HandlerFunc(ServeAsset)))
}
