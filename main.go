package main

import (
	"bytes"
	"flag"
	"fmt"

	"github.com/Tinker-Ware/digital-ocean-service/infrastructure"
	"github.com/Tinker-Ware/digital-ocean-service/interfaces"
	"github.com/Tinker-Ware/digital-ocean-service/usecases"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const defaultPath = "/etc/digital-ocean-service.conf"

var confFilePath = flag.String("conf", defaultPath, "Custom path for configuration file")

func main() {

	flag.Parse()

	config, err := infrastructure.GetConfiguration(*confFilePath)
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
		APIHost:     config.APIHost,
	}

	headers := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Authorization"})
	origins := handlers.AllowedOrigins([]string{"http://localhost", "http://192.168.33.10", "http://*.tinkerware.io", "https://*.tinkerware.io"})

	r := mux.NewRouter()

	subrouter := r.PathPrefix("/api/v1/cloud").Subrouter()

	subrouter.HandleFunc("/", handler.Login)
	subrouter.HandleFunc("/do_callback", handler.DOCallback).Methods("POST")
	subrouter.HandleFunc("/keys", handler.ShowKeys).Methods("GET")
	subrouter.HandleFunc("/keys", handler.CreateKey).Methods("POST")
	subrouter.HandleFunc("/instances", handler.CreateDroplet).Methods("POST")
	subrouter.HandleFunc("/instances", handler.ListDroplets).Methods("GET")
	subrouter.HandleFunc("/instance/{instanceID}", handler.GetInstance).Methods("GET")

	n := negroni.Classic()
	n.UseHandler(handlers.CORS(headers, origins)(r))

	port := bytes.Buffer{}

	port.WriteString(":")
	port.WriteString(config.Port)

	n.Run(port.String())

}
