package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Route-E-106/Frogfoot/server"
	_ "modernc.org/sqlite"
)

func main() {

	db, err := sql.Open("sqlite", "frogfoot.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	srv := server.NewServer(db)

	httpSrv := http.Server{
		Addr:    ":8080",
		Handler: srv.Routes(),
	}

	srv.Logger.Info("Starting server")
	err = httpSrv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
