package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stationapi/station-image-bucket/db"
	"github.com/stationapi/station-image-bucket/routes"
)

func main() {
	db, err := db.Connect()

	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)

	r.Post("/api/image/create", func(w http.ResponseWriter, r *http.Request) {
		err := routes.Create(w, r, db)

		if err != nil {
			fmt.Println(err)
		}
	})

	http.ListenAndServe(":3000", r)
}
