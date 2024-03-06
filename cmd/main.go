package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/database"
	"github.com/Noblefel/Rest-Api-Managemen-Kontak/internal/router"
)

const port = 8080

func main() {
	var (
		dbHost = flag.String("host", "localhost", "the database host")
		dbPort = flag.Int("port", 5432, "the database port")
		dbName = flag.String("name", "managemen_kontak", "the database name")
		dbUser = flag.String("u", "postgres", "the database user")
		dbPW   = flag.String("pw", "", "the database password")
	)
	flag.Parse()

	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s",
		*dbHost,
		*dbPort,
		*dbName,
		*dbUser,
		*dbPW,
	)

	db, err := database.Connect("pgx", dsn)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	log.Println("Starting server at port:", port)

	router := router.NewRouter(db)

	server := &http.Server{
		Addr:    fmt.Sprint("localhost:", port),
		Handler: router.Routes(),
	}

	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
