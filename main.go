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
The purpose of this service to simply increment a counter and display it to the caller, alongside an ID randomly
generated upon instantiation (but static throughout service lifetime), allowing for the identifying of the specific service
(e.g., if multiple instances are launched in parallel).

Source: https://github.com/chrfrantz/SimpleCountingService
*/

// ID of the service - identifying the instance. If not manually set, it will be populated with (reasonably) unique random numeric identifier.
var ID = "0"

// Background color for visual distinction of service's HTML output
var color = "#FFFFFF"

// Counter to be modified at runtime
var count = 0

// Key for http header to encode unique id for service (to simplify automated client-side analysis)
const headerKey = "Counter-ID"

/*
Adds service ID of the instance to the http header (Key: #headerKey). The purpose is to facilitate efficient client-side
identification of service ID for analytical purposes.
*/
func addCounterID(w http.ResponseWriter) {
	w.Header().Add(headerKey, ID)
}

/*
Generates style attribute for HTML output that produces unique background color for service instance (for simple
visual distinction of different service instances).
*/
func getHTMLStyle() string {
	return " style=\"background-color:#" + color + ";\""
}

/*
Handler to remind user to use the other handler (duplicate calls).
*/
func handlerRedirect(w http.ResponseWriter, r *http.Request) {

	// Add unique ID to http header
	addCounterID(w)

	// Prepare HTML response
	response := "<html" + getHTMLStyle() + "><head><title>Simple Service</title></head>" +
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

	// Add unique ID to http header
	addCounterID(w)

	// Increment counter
	count++

	// Prepare HTML response
	response := "<html" + getHTMLStyle() + "><head><title>Simple Counting Service</title></head>" +
		"<body><h1>Call to service (ID: " + ID +
		"); total calls: " + strconv.Itoa(count) + "</h1></body></html>"

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

	// Add unique ID to http header
	addCounterID(w)

	// Reset counter
	count = 0

	// Prepare HTML response
	response := "<html" + getHTMLStyle() + "><head><title>Simple Counting Service</title></head>" +
		"<body><h1>Call to service (ID: " + ID +
		"); counter reset!</h1></body></html>"

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

	// Set (reasonably) unique Identifier for UI output and background color by taking beginning of hashed timestamp
	if ID == "0" {
		h := fnv.New64a()
		_, err := h.Write([]byte(time.Now().String()))
		if err != nil {
			log.Fatal("Error during unique ID generation. Error: ", err.Error())
		}

		// Generate color (for service background) and trim for hex output (simply comment this code to suppress varying background colors)
		color = strconv.FormatUint(h.Sum64(), 16)[:6]

		// Generate reasonably unique (but readable) final ID
		ID = strconv.Itoa(int(h.Sum64()))[:3]
	}

	// Register endpoint handlers
	http.HandleFunc("/", handlerRedirect)
	http.HandleFunc("/count", handlerIncrement)
	http.HandleFunc("/reset", handlerReset)
	http.HandleFunc("/kill", handlerExitFailure)
	http.HandleFunc("/exit", handlerExitProper)

	// Launch service
	log.Println("Launching service on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
