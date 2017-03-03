package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"

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
	origins := handlers.AllowedOrigins([]string{"http://localhost", "http://provision.tinkerware.io", "https://provision.tinkerware.io"})

	r := mux.NewRouter()

	subrouter := r.PathPrefix("/api/v1/cloud").Subrouter()

	subrouter.HandleFunc("/digital_ocean/", handler.Login)
	subrouter.HandleFunc("/digital_ocean/oauth", handler.DOCallback).Methods("POST")
	subrouter.HandleFunc("/digital_ocean/keys", handler.ShowKeys).Methods("GET")
	subrouter.Handle("/digital_ocean/keys", interfaces.Adapt(http.HandlerFunc(handler.CreateKey), interfaces.GetToken(config.APIHost, config.Salt))).Methods("POST")
	subrouter.HandleFunc("/digital_ocean/instances", handler.CreateDroplet).Methods("POST")
	subrouter.HandleFunc("/digital_ocean/instances", handler.ListDroplets).Methods("GET")
	subrouter.HandleFunc("/digital_ocean/instance/{instanceID}", handler.GetInstance).Methods("GET")

	n := negroni.Classic()
	n.UseHandler(handlers.CORS(headers, origins)(r))

	port := bytes.Buffer{}

	port.WriteString(":")
	port.WriteString(config.Port)

	n.Run(port.String())

}
