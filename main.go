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

var id = 0
var count = 0

/*
Handler to remind user to use the other handler (duplicate calls).
 */
func handlerRedirect(w http.ResponseWriter, r *http.Request) {

	response := "<html><head><title>Simple Service</title></head>" +
		"<body><h1>Please call service on path '/count'.</h1></body></html>"
	fmt.Fprintln(w, response)

}

/*
Returns a static identifier, as well as a incrementing count.
 */
func handlerDefault(w http.ResponseWriter, r *http.Request) {

	response := "<html><head><title>Simple Service</title></head>" +
		"<body><h1>Call to service " + strconv.Itoa(id) +
		"; total calls: " + strconv.Itoa(count) + "</h1></body></html>"

	fmt.Fprintln(w, response)

	count++

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
	http.HandleFunc("/count", handlerDefault)
	log.Println("Launching service on port " + port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
