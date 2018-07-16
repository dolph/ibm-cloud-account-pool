package main

import "encoding/json"
import "net/http"

import "github.com/gorilla/mux"

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
