package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Kavan03/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

//Defines a struct apiConfig that holds a reference to the database queries object.
//apiConfig is a dependency injection struct, making it easier to pass database queries (database.Queries) into API handlers without needing global variables.
type apiConfig struct{
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port not found")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Database URL not found")
	}

	//connect database
	conn, err := sql.Open("postgres",dbURL)
	if err!=nil {
		log.Fatal("Can't connect with Database:",err)
	}
	
	dbQueries := database.New(conn)
	//In a typical web application, you'll create an instance of apiConfig at the start of your program (often in the main function). This instance is then passed to HTTP handlers or other parts of the application that need to access the database.
	apiCfg := apiConfig{
		DB: dbQueries,
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
	v1Router.Post("/users",apiCfg.handlerCreateUser)
	v1Router.Get("/users",apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds",apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds",apiCfg.handlerGetFeeds)
	v1Router.Post("/feed_follow",apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follow",apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollow))
	v1Router.Delete("/feed_follow/{feed_follow_id}",apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerPostsGet))
	router.Mount("/v1",v1Router)

	srv := &http.Server{ //Configuring the server using build in package http.Server
		Handler: router, //Sets chi router as request handler
		Addr: ":" + portString, //The server listens on the specified port (portString).
	}

	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)

	log.Printf("Server starting on Port %v",portString)
	err = srv.ListenAndServe() // Starts the HTTP server, listens for incoming requests.
	if err != nil {
		log.Fatal(err)
	}
}