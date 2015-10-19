// Copyright 2015 The Loadcat Authors. All rights reserved.

package ui

import (
	"net/http"
)

func init() {
	Router.NewRoute().
		Methods("GET").
		PathPrefix("/assets").
		Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./ui/assets"))))
}
