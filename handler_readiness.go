package main

import "net/http"

//This function is an HTTP handler for checking the readiness of a service.

func handlerReadiness(w http.ResponseWriter, r *http.Request){
	respondWithJSON(w,200,struct{}{})
}