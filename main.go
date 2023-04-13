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
(e.g., if multiple instances are launched in parallel). To support the identification of service instances, it returns a
static background color for endpoints with visual output and passes count information in the http header (which is used by
the complementary client that counts invocation distribution across unique instances of the service).

Source: https://github.com/chrfrantz/SimpleCountingService
*/

// ID of the service - identifying the instance. If not manually set, it will be populated with (reasonably) unique
// random numeric identifier and remains static throughout the service runtime.
var ID = "0"

// Background color for visual distinction of service's HTML output
var color = "#FFFFFF"

// Counter to be modified at runtime
var count = 0

// Key for http header to encode unique id for service (to simplify automated client-side analysis)
const HeaderKey = "Counter-ID"

// Establish compatibility with PaaS platforms
const EnvVarPort = "PORT"

// Default port
const DefaultPort = "8080"

// Switch to control printing to console (e.g., for performance reasons)
const PrintToConsole = false

/*
Service endpoint constants
*/

const PathCount = "/count"

const PathReset = "/reset"

const PathKill = "/kill"

const PathExit = "/exit"

/*
Function to override printing depending on activation of #PrintToConsole.
*/
func Println(content string) {
	if PrintToConsole {
		fmt.Println(content)
	}
}

/*
Adds service ID of the instance to the http header (Key: #headerKey). The purpose is to facilitate efficient client-side
identification of service ID for analytical purposes.
*/
func addCounterID(w http.ResponseWriter) {
	w.Header().Add(HeaderKey, ID)
}

/*
Adds HTML content type to HTTP response.
*/
func addHTMLContentType(w http.ResponseWriter) {
	w.Header().Add("content-type", "text/html")
}

/*
Generates style attribute for HTML output that produces unique background color for service instance (for simple
visual distinction of different service instances).
*/
func getHTMLStyle() string {
	return " style=\"background-color:#" + color + ";\""
}

/*
Returns complete HTML output to be used as a response to service request.
Takes content of body (message) as input and wraps it with standard HTML.
*/
func generateHTMLOutput(content string) string {
	return "<html" + getHTMLStyle() + "><head><title>Simple Service</title></head><body><h1>" +
		content + "</h1></body></html>"
}

/*
Handler to remind user to use the other handler (duplicate calls).
*/
func handlerRedirect(w http.ResponseWriter, r *http.Request) {

	// Sets content type for html output
	addHTMLContentType(w)

	// Add unique ID to http header
	addCounterID(w)

	// Prepare HTML response
	response := generateHTMLOutput("The SimpleCounterService serves the purpose of illustrating the repeated " +
		"invocation of web services (e.g., to demonstrate load balancing strategies).<br><br>" +
		"It has the following API:<br>" +
		" - Path '<a href=\"" + PathCount + "\">" + PathCount +
		"</a>' increments the count for service invocations.<br>" +
		" - Path '<a href=\"" + PathReset + "\">" + PathReset + "</a>' resets the count.<br>" +
		" - To explore exit behaviour (and potential recovery), use path '<a href=\"" + PathKill + "\">" + PathKill +
		"</a>' (killing service; with error) and '<a href=\"" + PathExit + "\">" + PathExit +
		"</a>' (proper service termination without error).")

	// Sleep for one second
	time.Sleep(1000)

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

	// Sets content type for html output
	addHTMLContentType(w)

	// Add unique ID to http header
	addCounterID(w)

	// Increment counter
	count++

	// Print to console
	Println("Counter " + ID + ": Call on path " + PathCount + ": " + strconv.Itoa(count) + " invocation(s)")

	// Prepare HTML response
	response := generateHTMLOutput("Call to service (ID: " + ID + "); total calls: " + strconv.Itoa(count))

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

	// Sets content type for html output
	addHTMLContentType(w)

	// Add unique ID to http header
	addCounterID(w)

	// Reset counter
	count = 0

	// Print to console
	Println("Counter " + ID + ": Call on path " + PathReset + ".")

	// Prepare HTML response
	response := generateHTMLOutput("Call to service (ID: " + ID + "); counter reset!")

	// Return response
	_, err := fmt.Fprintln(w, response)
	if err != nil {
		http.Error(w, "Error when writing response. Error: "+err.Error(), http.StatusInternalServerError)
	}
}

/*
Kills service with error - just to demonstrate service recovery (e.g., via restart policies)
*/
func handlerExitFailure(w http.ResponseWriter, r *http.Request) {

	// Print to console
	Println("Counter " + ID + ": Call on path " + PathKill + ".")

	// Exit service with error
	os.Exit(1)
}

/*
Kills service without error - just to demonstrate service recovery (e.g., via restart policies)
*/
func handlerExitProper(w http.ResponseWriter, r *http.Request) {

	// Print to console
	Println("Counter " + ID + ": Call on path " + PathExit + ".")

	// Exit service properly
	os.Exit(0)
}

func main() {

	// Check for custom port
	port := os.Getenv(EnvVarPort)
	if port == "" {
		port = DefaultPort
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
	http.HandleFunc(PathCount, handlerIncrement)
	http.HandleFunc(PathReset, handlerReset)
	http.HandleFunc(PathKill, handlerExitFailure)
	http.HandleFunc(PathExit, handlerExitProper)

	// Launch service
	log.Println("Launching SimpleCountingService (with unique ID " + ID + ") on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
