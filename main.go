package main

import "fmt"
import "io/ioutil"
import "log"
import "math/rand"
import "net/http"
import "os"
import "time"

import "github.com/gorilla/mux"
import "gopkg.in/yaml.v2"

func main() {
	// Seed random number generator for token & ID generation.
	rand.Seed(time.Now().UTC().UnixNano())

	// Get configuration from the environment.
	iam_apikey := os.Getenv("CLOUDANT_IAM_APIKEY")
	if iam_apikey == "" {
		fmt.Println("Missing an API key to access Cloudant (CouchDB): CLOUDANT_IAM_APIKEY")
		os.Exit(1)
	}
	couchdb_url := os.Getenv("CLOUDANT_IAM_URL")
	if couchdb_url == "" {
		fmt.Println("Missing an endpoint for Cloudant (CouchDB): CLOUDANT_IAM_URL")
		os.Exit(1)
	}

	auth := NewIAM(iam_apikey)
	dbclient := NewClient(auth, couchdb_url)

	// Test database connectivity.
	if _, err := dbclient.ListAllDatabases(); err == nil {
		fmt.Println("Database connected succesfully.")
	} else {
		fmt.Println("Unable to connect to database")
		panic(err)
	}

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
