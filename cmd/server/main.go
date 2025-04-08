package main

import (
	"log"
	"net/http"

	"github.com/Route-E-106/Frogfoot/server"
)

func main() {

	srv := server.NewServer()

	httpSrv := http.Server{
		Addr:    ":8080",
		Handler: srv.Routes(),
	}

	err := httpSrv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
