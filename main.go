package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

/*
The purpose of this service to simply increment a counter and display it to the caller, alongside a static id
identifying the service instance.
 */

// ID of the service - identifying the instance
const ID = 0
// Counter to be modified at runtime
var count = 0

/*
Handler to remind user to use the other handler (duplicate calls).
 */
func handlerRedirect(w http.ResponseWriter, r *http.Request) {

	// Prepare HTML response
	response := "<html><head><title>Simple Service</title></head>" +
		"<body><h1>Please call service on path '/count'.</h1></body></html>"

	// Return response
	fmt.Fprintln(w, response)
}

/*
Returns a static identifier, as well as a incrementing count.
 */
func handlerIncrement(w http.ResponseWriter, r *http.Request) {

	// Increment counter
	count++

	// Prepare HTML response
	response := "<html><head><title>Simple Counting Service</title></head>" +
		"<body><h1>Call to service " + strconv.Itoa(ID) +
		"; total calls: " + strconv.Itoa(count) + "</h1></body></html>"

	// Return response
	fmt.Fprintln(w, response)
}

/*
Resets the counter
 */
func handlerReset(w http.ResponseWriter, r *http.Request) {

	// Reset counter
	count = 0

	// Prepare HTML response
	response := "<html><head><title>Simple Counting Service</title></head>" +
		"<body><h1>Call to service " + strconv.Itoa(ID) +
		"; counter reset!</h1></body></html>"

	// Return response
	fmt.Fprintln(w, response)
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

	http.HandleFunc("/", handlerRedirect)
	http.HandleFunc("/count", handlerIncrement)
	http.HandleFunc("/reset", handlerReset)
	log.Println("Launching service on port " + port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
