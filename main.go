package main

import "fmt"
import "math/rand"
import "log"
import "time"
import "net/http"

import "github.com/gorilla/mux"
import "github.com/ernsheong/grand"

func main() {
	// Seed random number generator for token & ID generation.
	rand.Seed(time.Now().UTC().UnixNano())

	token := NewToken()
	fmt.Println("Initial token:", token.Id)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/reservations", CreateReservation).Methods("POST")
	router.HandleFunc("/reservations/{reservationId}", GetReservation).Methods("GET")
	router.HandleFunc("/reservations/{reservationId}", DeleteReservation).Methods("DELETE")

	fmt.Println("Listening on localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func NewToken() Token {
	generator := grand.NewGenerator(grand.CharSetBase62)
	token := Token{Id: generator.GenerateRandomString(8)}
	tokens[token.Id] = &token
	return token
}
