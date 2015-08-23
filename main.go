package main

import (
	"bytes"
	"flag"
	"fmt"

	"github.com/codegangsta/negroni"
	"github.com/digital-ocean-service/infraestructure"
	"github.com/digital-ocean-service/interfaces"
	"github.com/digital-ocean-service/usecases"
	"github.com/gorilla/mux"
)

const defaultPath = "/etc/digital-ocean-service.conf"

var confFilePath = flag.String("conf", defaultPath, "Custom path for configuration file")

func main() {

	flag.Parse()

	config, err := infraestructure.GetConfiguration(*confFilePath)
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot parse configuration")
	}

	doInteractor := usecases.DOInteractor{}

	handler := interfaces.WebServiceHandler{
		Interactor:  doInteractor,
		ID:          config.ClientID,
		Secret:      config.ClientSecret,
		Scopes:      config.Scopes,
		RedirectURI: config.RedirectURI,
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

	port := bytes.Buffer{}

	port.WriteString(":")
	port.WriteString(config.Port)

	n.Run(port.String())

}
