package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string){
	if code>499{
		log.Println("Responding with 5XX error:",msg)
	}
	type errResponse struct{
		Error string `json:"error"`
	}
	respondWithJSON(w,code,errResponse{
		Error: msg,
	})
}

// w http.ResponseWriter → Used to send the response back to the client.
// code int → Expected HTTP status code (e.g., 200, 400, 500).
// payload interface{} → Any Go data structure (struct, map, slice) that should be converted into JSON.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){
	dat, err := json.Marshal(payload) // Convert the payload into JSON byte array
	if err!=nil{
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type","application/json")
	w.WriteHeader(code)
	w.Write(dat)
}