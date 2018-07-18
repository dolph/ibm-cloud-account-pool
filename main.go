package main

import "fmt"
import "math/rand"
import "log"
import "io/ioutil"
import "time"
import "net/http"

import "github.com/gorilla/mux"
import "gopkg.in/yaml.v2"

func main() {
	// Seed random number generator for token & ID generation.
	rand.Seed(time.Now().UTC().UnixNano())

	data := loadData()
	for _, token := range data.Tokens {
		tokens[token.Id] = &token
	}
	fmt.Printf("Loaded %v token(s).\n", len(tokens))
	for _, account := range data.Accounts {
		accounts[account.Username] = &account
	}
	fmt.Printf("Loaded %v account(s).\n", len(accounts))

	fmt.Println("Listening on localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", NewRouter()))
}

func NewRouter() http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/reservations", CreateReservation).Methods("POST").Queries("token", "{[a-zA-Z0-9]8}")
	router.HandleFunc("/reservations/{reservationId}", GetReservation).Methods("GET").Queries("token", "{[a-zA-Z0-9]8}")
	router.HandleFunc("/reservations/{reservationId}", DeleteReservation).Methods("DELETE").Queries("token", "{[a-zA-Z0-9]8}")
	return router
}

func loadData() Yaml {
	var data Yaml
	yamlFile, err := ioutil.ReadFile("data.yml")
	if err != nil {
		log.Fatalf("Error loading data.yml: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Fatalf("Error unmarshalling data.yml: %v", err)
	}
	return data
}

type Yaml struct {
	Tokens   []Token
	Accounts []Account
}
