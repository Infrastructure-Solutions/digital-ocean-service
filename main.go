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
		RedirectURI: "http://localhost:7000/do_callback",
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler.Login)

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":7000")

}
