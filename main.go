package main

import (
	"github.com/codegangsta/negroni"
	"github.com/digital-ocean-service/interfaces"
	"github.com/digital-ocean-service/usecases"
	"github.com/gorilla/mux"
)

func main() {
	doInteractor := usecases.DOInteractor{}

	handler := interfaces.WebServiceHandler{
		Interactor:  doInteractor,
		ID:          "",
		Secret:      "",
		RedirectURI: "http://localhost:7000/do_callback",
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler.Login)
	r.HandleFunc("/do_callback", handler.DOCallback).Methods("GET")
	r.HandleFunc("/keys", handler.ShowKeys).Methods("GET")
	r.HandleFunc("/keys", handler.CreateKey).Methods("POST")
	r.HandleFunc("/droplets", handler.CreateDroplet).Methods("POST")
	r.HandleFunc("/droplets", handler.ListDroplets).Methods("GET")

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":7000")

}
