package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Kavan03/rssagg/internal/database"
	"github.com/google/uuid"
)

//This function is an HTTP handler for checking the readiness of a service.

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Error parsing JSON: %v",err))
		return
	}
	user, err:= apiCfg.DB.CreateUser(r.Context(),database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
	})
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Couldn't create user: %v",err))
		return
	}
	respondWithJSON(w,201,databaseUserToUser(user))
}

// This Go function handlerCreateUser is part of an HTTP handler that processes an incoming HTTP request, extracts user data from the request body (in JSON format), and creates a user in the database.
// The function is a method of the apiConfig struct. The method takes two arguments:
// w http.ResponseWriter: This is used to write the HTTP response back to the client.
// r *http.Request: This is the incoming HTTP request, containing details like request parameters, body, etc.
// parameters Struct: A local struct parameters is defined to represent the expected data format in the request body. Here, it expects a name field of type string (deserialized from JSON using the tag json:"name").
// A JSON decoder (decoder) is created to read the request body (r.Body).
// The Decode function attempts to parse the incoming JSON into the params variable (which is of type parameters). If successful, params.Name will be populated with the data from the request.
// If an error occurs while decoding the JSON (e.g., invalid JSON format or missing required fields), the handler responds with an HTTP 400 (Bad Request) status code and a descriptive error message.
// The apiCfg.DB.CreateUser function is called, which likely interacts with the database to create a new user record.
// The CreateUserParams struct is populated with:
// ID: A newly generated unique user ID using uuid.New().
// CreatedAt and UpdatedAt: The current UTC time (time.Now().UTC()).
// Name: The name of the user, which was parsed from the incoming JSON (params.Name).
// If the user creation fails (e.g., due to a database issue), the handler responds with an HTTP 400 status code and an error message.
// If the user is successfully created, the handler responds with HTTP status code 200 (OK) and the user object (which presumably contains the created user data) in JSON format.

// This handler:
// Receives a POST request to create a new user.
// Extracts the user's name from the JSON body of the request.
// Creates a new user in the database with a generated UUID, the current UTC time, and the provided name.
// Handles errors for invalid JSON input or failure to create the user.
// Responds with the created user if successful, or an error message if something goes wrong.

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User){
	respondWithJSON(w,200,databaseUserToUser(user))
}