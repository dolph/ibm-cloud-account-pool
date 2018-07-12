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
	router.HandleFunc("/reservations", Reservations).Methods("POST")
	router.HandleFunc("/reservations/{reservationId}", GetReservation).Methods("GET")
	router.HandleFunc("/reservations/{reservationId}", DeleteReservation).Methods("DELETE")
	fmt.Println("Listening on localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func Reservations(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func CreateReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["token"]
	fmt.Fprintf(w, "Token ID: %q", html.EscapeString(tokenId))

	reservation := Reservation{}
	json.NewEncoder(w).Encode(reservation)
}

func GetReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reservationId := vars["reservationId"]
	fmt.Fprintf(w, "Reservation ID: %q", html.EscapeString(reservationId))
	reservation := Reservation{
		Id: reservationId,
	}
	json.NewEncoder(w).Encode(reservation)
}

func DeleteReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reservationId := vars["reservationId"]
	fmt.Fprintf(w, "Delete Reservation ID: %q", html.EscapeString(reservationId))
}
