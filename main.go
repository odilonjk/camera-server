package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/record/start", startRecord).Methods("POST")
	router.HandleFunc("/api/v1/record/stop", stopRecord).Methods("POST")
	router.HandleFunc("/api/v1/record", getRecord).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func startRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
}

func stopRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
}

func getRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
}
