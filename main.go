package main

import (
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

/*
The purpose of this service to simply increment a counter and display it to the caller, alongside a static id
identifying the service instance.
*/

// ID of the service - identifying the instance. If not manually set, it will be populated with unique hash.
var ID = "0"

// Counter to be modified at runtime
var count = 0

/*
Handler to remind user to use the other handler (duplicate calls).
*/
func handlerRedirect(w http.ResponseWriter, r *http.Request) {

	// Prepare HTML response
	response := "<html><head><title>Simple Service</title></head>" +
		"<body><h1>Please call service on path '/count'. " +
		"To explore exit behaviour, use path '/kill' (error) and '/exit' (no error) </h1></body></html>"

	// Return response
	_, err := fmt.Fprintln(w, response)
	if err != nil {
		http.Error(w, "Error when writing response. Error: "+err.Error(), http.StatusInternalServerError)
	}
}

/*
Returns a static identifier, as well as a incrementing count.
*/
func handlerIncrement(w http.ResponseWriter, r *http.Request) {

	// Increment counter
	count++

	// Prepare HTML response
	response := "<html><head><title>Simple Counting Service</title></head>" +
		"<body><h1>Call to service " + ID +
		"; total calls: " + strconv.Itoa(count) + "</h1></body></html>"

	// Return response
	_, err := fmt.Fprintln(w, response)
	if err != nil {
		http.Error(w, "Error when writing response. Error: "+err.Error(), http.StatusInternalServerError)
	}
}

/*
Resets the counter
*/
func handlerReset(w http.ResponseWriter, r *http.Request) {

	// Reset counter
	count = 0

	// Prepare HTML response
	response := "<html><head><title>Simple Counting Service</title></head>" +
		"<body><h1>Call to service " + ID +
		"; counter reset!</h1></body></html>"

	// Return response
	_, err := fmt.Fprintln(w, response)
	if err != nil {
		http.Error(w, "Error when writing response. Error: "+err.Error(), http.StatusInternalServerError)
	}
}

/*
Kills service with error - just to demonstrate restart policies
*/
func handlerExitFailure(w http.ResponseWriter, r *http.Request) {

	// Exit service with error
	os.Exit(1)
}

/*
Kills service without error - just to demonstrate restart policies
*/
func handlerExitProper(w http.ResponseWriter, r *http.Request) {

	// Exit service properly
	os.Exit(0)
}

// Establish compatibility with PaaS platforms
const ENV_VAR_PORT = "PORT"

// Default port
const DEFAULT_PORT = "8080"

func main() {

	// Check for custom port
	port := os.Getenv(ENV_VAR_PORT)
	if port == "" {
		port = DEFAULT_PORT
	}

	// Set (reasonably) unique Identifier for UI by taking beginning of hashed timestamp
	if ID == "0" {
		h := fnv.New32a()
		_, err := h.Write([]byte(time.Now().String()))
		if err != nil {
			log.Fatal("Error during unique ID generation. Error: ", err.Error())
		}
		// Generate final ID and remove suffix
		ID = strconv.Itoa(int(h.Sum32()))[:8]
	}

	http.HandleFunc("/", handlerRedirect)
	http.HandleFunc("/count", handlerIncrement)
	http.HandleFunc("/reset", handlerReset)
	http.HandleFunc("/kill", handlerExitFailure)
	http.HandleFunc("/exit", handlerExitProper)
	log.Println("Launching service on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
