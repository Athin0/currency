package main

import (
	"currency/inretnal/adapters"
	"currency/inretnal/ports"
	"currency/inretnal/service"
	"log"
	"net/http"
)

func main() {
	pg, err := adapters.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	app := service.NewApp(pg)
	server := ports.NewHTTPServer("8080", app)

	err = http.ListenAndServe(":8080", server.Handler())
	if err != nil {
		log.Fatal(err)
	}
}
