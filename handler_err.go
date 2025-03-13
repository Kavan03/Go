package main

import "net/http"

//This function is an HTTP handler for checking the readiness of a service.

func handlerErr(w http.ResponseWriter, r *http.Request){
	respondWithError(w,400,"Something went wrong")
}