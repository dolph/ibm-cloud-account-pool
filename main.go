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

func Index(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(statistics)
}

func CreateReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tokenId := vars["token"]
	token := Token{
		Id: tokenId,
	}
	reservation := Reservation{
		Token: &token,
	}

	// TODO: Add reservation to reservations

	json.NewEncoder(w).Encode(reservation)
}

func GetReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reservation := reservations[vars["reservationId"]]
	json.NewEncoder(w).Encode()
}

func DeleteReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reservationId := vars["reservationId"]
	reservation := reservations[reservationId]

	// TODO: Release reservation

	json.NewEncoder(w).Encode(reservation)
}
