package main

import (
	"log"
	"net/http"

	"github.com/maxkaiser11/wowcombatparser/internal/api"
)

func main() {
	srv, err := api.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", srv.Routes()); err != nil {
		log.Fatal(err)
	}
}
