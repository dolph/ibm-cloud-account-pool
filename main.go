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

	config := loadConfig()
	for _, token := range config.Tokens {
		tokens[token.Id] = &token
	}
	fmt.Printf("Loaded %v token(s).\n", len(tokens))
	for _, account := range config.Accounts {
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

func loadConfig() Yaml {
	var config Yaml
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Error loading config.yml: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling config.yml: %v", err)
	}
	return config
}

type Yaml struct {
	Tokens   []Token
	Accounts []Account
}
