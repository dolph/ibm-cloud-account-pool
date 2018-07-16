package main

import "encoding/json"
import "fmt"
import "html"
import "log"
import "net/http"

import "github.com/gorilla/mux"

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/reservations", CreateReservation).Methods("POST")
	router.HandleFunc("/reservations/{reservationId}", GetReservation).Methods("GET")
	router.HandleFunc("/reservations/{reservationId}", DeleteReservation).Methods("DELETE")
	fmt.Println("Listening on localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
