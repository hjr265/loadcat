// Copyright 2015 The Loadcat Authors. All rights reserved.

package ui

import (
	"net/http"
)

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/balancers", http.StatusSeeOther)
}

func init() {
	Router.NewRoute().
		Methods("GET").
		Path("/").
		Handler(http.HandlerFunc(ServeIndex))
}
