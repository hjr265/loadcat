// Copyright 2015 The Loadcat Authors. All rights reserved.

package ui

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"gopkg.in/mgo.v2/bson"

	"github.com/hjr265/loadcat/data"
	"github.com/hjr265/loadcat/feline"
)

func ServeBalancerList(w http.ResponseWriter, r *http.Request) {
	bals, err := data.ListBalancers()
	if err != nil {
		panic(err)
	}

	err = TplBalancerList.Execute(w, struct {
		Balancers []data.Balancer
	}{
		Balancers: bals,
	})
	if err != nil {
		panic(err)
	}
}

func ServeBalancerNewForm(w http.ResponseWriter, r *http.Request) {
	err := TplBalancerNewForm.Execute(w, struct {
	}{})
	if err != nil {
		panic(err)
	}
}

func HandleBalancerCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	body := struct {
		Label string `schema:"label"`
	}{}
	err = schema.NewDecoder().Decode(&body, r.PostForm)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	bal := data.Balancer{}
	bal.Label = body.Label
	bal.Settings.Hostname = "example.com"
	bal.Settings.Port = 80
	err = bal.Put()
	if err != nil {
		panic(err)
	}

	err = feline.Commit(&bal)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/balancers/"+bal.Id.Hex()+"/edit", http.StatusSeeOther)
}

func ServeBalancer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if !bson.IsObjectIdHex(vars["id"]) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	bal, err := data.GetBalancer(bson.ObjectIdHex(vars["id"]))
	if err != nil {
		panic(err)
	}

	err = TplBalancerView.Execute(w, struct {
		Balancer *data.Balancer
	}{
		Balancer: bal,
	})
	if err != nil {
		panic(err)
	}
}

func ServeBalancerEditForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if !bson.IsObjectIdHex(vars["id"]) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	bal, err := data.GetBalancer(bson.ObjectIdHex(vars["id"]))
	if err != nil {
		panic(err)
	}

	err = TplBalancerEditForm.Execute(w, struct {
		Balancer     *data.Balancer
		Protocols    []data.Protocol
		SSLVerifys   []data.SSLVerify
		Algorithms   []data.Algorithm
		CipherSuites []data.CipherSuite
	}{
		Balancer:     bal,
		Protocols:    data.Protocols,
		SSLVerifys:   data.SSLVerifys,
		Algorithms:   data.Algorithms,
		CipherSuites: data.CipherSuites,
	})
	if err != nil {
		panic(err)
	}
}

func HandleBalancerUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if !bson.IsObjectIdHex(vars["id"]) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	bal, err := data.GetBalancer(bson.ObjectIdHex(vars["id"]))
	if err != nil {
		panic(err)
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	body := struct {
		Label    string `schema:"label"`
		Settings struct {
			Hostname   string `schema:"hostname"`
			Port       int    `schema:"port"`
			Protocol   string `schema:"protocol"`
			Algorithm  string `schema:"algorithm"`
			SSLOptions struct {
				CipherSuite string  `schema:"cipher_suite"`
				Certificate *string `schema:"certificate"`
				PrivateKey  *string `schema:"private_key"`
				SSLVerify    string `schema:"sslverify"`
			} `schema:"ssl_options"`
			SSLVerifyClient struct {
				ClientCertificate *string `schema:"clientcertificate"`
			}
		} `schema:"settings"`
	}{}
	err = schema.NewDecoder().Decode(&body, r.PostForm)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	bal.Label = body.Label
	bal.Settings.Hostname = body.Settings.Hostname
	bal.Settings.Port = body.Settings.Port
	bal.Settings.Protocol = data.Protocol(body.Settings.Protocol)
	bal.Settings.Algorithm = data.Algorithm(body.Settings.Algorithm)
	if body.Settings.Protocol == "https" {
		bal.Settings.SSLOptions.CipherSuite = "recommended"
		bal.Settings.SSLOptions.SSLVerify = "off"
		if body.Settings.SSLOptions.SSLVerify == "on" {
			bal.Settings.SSLOptions.SSLVerify = data.SSLVerify(body.Settings.SSLOptions.SSLVerify)
			if body.Settings.SSLVerifyClient.ClientCertificate != nil {
				bal.Settings.SSLVerifyClient.ClientCertificate = []byte(*body.Settings.SSLVerifyClient.ClientCertificate)
			}
		}
		if body.Settings.SSLOptions.Certificate != nil {
			bal.Settings.SSLOptions.Certificate = []byte(*body.Settings.SSLOptions.Certificate)
		}
		if body.Settings.SSLOptions.PrivateKey != nil {
			bal.Settings.SSLOptions.PrivateKey = []byte(*body.Settings.SSLOptions.PrivateKey)
		}
	}
	err = bal.Put()
	if err != nil {
		panic(err)
	}

	err = feline.Commit(bal)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/balancers/"+bal.Id.Hex(), http.StatusSeeOther)
}

func HandleBalancerDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if !bson.IsObjectIdHex(vars["id"]) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	bal, err := data.GetBalancer(bson.ObjectIdHex(vars["id"]))
	if err != nil {
		panic(err)
	}

	err = bal.Delete()
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/balancers", http.StatusSeeOther)
}

func init() {
	Router.NewRoute().
		Methods("GET").
		Path("/balancers").
		Handler(http.HandlerFunc(ServeBalancerList))
	Router.NewRoute().
		Methods("GET").
		Path("/balancers/new").
		Handler(http.HandlerFunc(ServeBalancerNewForm))
	Router.NewRoute().
		Methods("POST").
		Path("/balancers/new").
		Handler(http.HandlerFunc(HandleBalancerCreate))
	Router.NewRoute().
		Methods("GET").
		Path("/balancers/{id}").
		Handler(http.HandlerFunc(ServeBalancer))
	Router.NewRoute().
		Methods("GET").
		Path("/balancers/{id}/edit").
		Handler(http.HandlerFunc(ServeBalancerEditForm))
	Router.NewRoute().
		Methods("POST").
		Path("/balancers/{id}/edit").
		Handler(http.HandlerFunc(HandleBalancerUpdate))
	Router.NewRoute().
		Methods("POST").
		Path("/balancers/{id}/delete").
		Handler(http.HandlerFunc(HandleBalancerDelete))
}
