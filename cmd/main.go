package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/database"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/router"
	"github.com/joho/godotenv"
)

const port = 8080

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	log.Println("Starting server at port:", port)

	router := router.NewRouter()

	server := &http.Server{
		Addr:    fmt.Sprint("localhost:", port),
		Handler: router.Routes(),
	}

	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
