package main

import "encoding/json"
import "net/http"

import "github.com/gorilla/mux"

func SendJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func Index(w http.ResponseWriter, r *http.Request) {
	SendJSON(w, statistics)
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

	SendJSON(w, reservation)
}

func GetReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reservation := reservations[vars["reservationId"]]
	SendJSON(w, reservation)
}

func DeleteReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reservationId := vars["reservationId"]
	reservation := reservations[reservationId]

	// TODO: Release reservation

	SendJSON(w, reservation)
}
