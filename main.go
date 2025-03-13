package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port not found")
	}

	router := chi.NewRouter() // Creating router

	router.Use(cors.Handler(cors.Options{ // Adding cors so that users can request from browser
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz",handlerReadiness)
	v1Router.Get("/err",handlerErr)
	router.Mount("/v1",v1Router)

	srv := &http.Server{ //Configuring the server using build in package http.Server
		Handler: router, //Sets chi router as request handler
		Addr: ":" + portString, //The server listens on the specified port (portString).
	}
	log.Printf("Server starting on Port %v",portString)
	err := srv.ListenAndServe() // Starts the HTTP server, listens for incoming requests.
	if err != nil {
		log.Fatal(err)
	}
}