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
		ID:          "f584247261a56d4003d795842fbeaacdd82d5624693bd6e00c02b3c6d675cf44",
		Secret:      "9f515e7485d3d8e602b726620fe115084ff95f641228285901bd1959a350c05c",
		RedirectURI: "http://localhost:7000/do_callback",
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler.Login)
	r.HandleFunc("/do_callback", handler.DOCallback).Methods("GET")
	r.HandleFunc("/keys", handler.ShowKeys).Methods("GET")
	r.HandleFunc("/keys", handler.CreateKey).Methods("POST")

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":7000")

}
