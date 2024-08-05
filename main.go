package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Noblefel/baic-rest-api-kontak/internal/handlers"
	"github.com/Noblefel/baic-rest-api-kontak/internal/storage"
	"github.com/Noblefel/baic-rest-api-kontak/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

const port = 8080

var appStorage = storage.New()
var appHandlers = handlers.New(appStorage)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	pw, _ := bcrypt.GenerateFromPassword([]byte("asdasdasd"), 10)
	appStorage.Register("asd@example.com", string(pw))

	server := &http.Server{
		Addr:    fmt.Sprint("localhost:", port),
		Handler: routes(),
	}

	log.Println("Starting server at port:", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func routes() http.Handler {
	api := http.NewServeMux()

	api.HandleFunc("POST /register", appHandlers.Register)
	api.HandleFunc("POST /login", appHandlers.Login)

	api.Handle("GET /me", auth(appHandlers.GetAccount))
	api.Handle("PUT /me", auth(appHandlers.UpdateAccount))
	api.Handle("DELETE /me", auth(appHandlers.DeleteAccount))
	api.Handle("GET /me/contacts", auth(appHandlers.GetAccountContacts))

	api.Handle("POST /contacts", auth(appHandlers.CreateContact))
	api.Handle("GET /contacts/{id}", auth(appHandlers.GetContact))
	api.Handle("PUT /contacts/{id}", auth(appHandlers.UpdateContact))
	api.Handle("DELETE /contacts/{id}", auth(appHandlers.DeleteContact))

	api.HandleFunc("GET /ping", handlers.Ping)

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", api))
	return mux
}

func auth(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			utils.Message(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		userId, err := utils.VerifyJWT(tokenString)
		if err != nil {
			utils.Message(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		_, err = appStorage.GetUser(userId)
		if err != nil {
			utils.Message(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
