package main

import (
	"fmt"
	"net/http"
	"strconv"
)

/*
The purpose of this client is to automate the repeated interrogation of the deployed SimpleCountingService and establish a
distribution of responses (including errors, as well as across different service instances in the case of load balancing).

Source: https://github.com/chrfrantz/SimpleCountingService
*/

/*
SimpleCountingService Endpoint of interest for analysis
*/
const URL = "http://localhost:8080/count"

/*
Number of invocations to be performed
*/
const Rounds = 50

/*
Http header key that holds unique id for service
*/
const HeaderKey = "Counter-ID"

/*
Constants to record calls to endpoint that failed
*/
const ErrorKey = "Error"

/*
Map to collect responses for later analysis
*/
var responseMap = make(map[string]int)

func main() {

	// Counter for performed runs
	ct := 0

	// Instantiate client
	client := &http.Client{}
	defer client.CloseIdleConnections()

	fmt.Println("Running " + strconv.Itoa(Rounds) + " invocations on service at " + URL)

	// Iterate through calls to specified endpoint
	for i := 0; i < Rounds; i++ {

		fmt.Println("Performing Round " + strconv.Itoa(i))

		// Perform request
		resp, err := client.Get(URL)
		if err != nil {
			fmt.Println("Error during request " + strconv.Itoa(ct) + ": " + err.Error())
			// Count as error
			addCountForService(ErrorKey)
		} else {
			// Retrieve identifier from http header
			serviceKey := resp.Header.Get(HeaderKey)
			// Count response
			addCountForService(serviceKey)
		}
	}

	// Analyze results

	fmt.Println("\nResults:\n--------")

	for k, v := range responseMap {
		// Calculate percentage relative to rounds
		percentage := fmt.Sprintf("%f", float64(v)/float64(Rounds)*100)
		// Print accordingly
		fmt.Println("Service " + k + ": " + percentage + " percent")
	}

}

/*
Handles the counting of responses for given service instances. Requires service identifier as parameter.
*/
func addCountForService(serviceKey string) {
	// Aggregate counts for service instance
	if serviceKey != "" {
		v, ok := responseMap[serviceKey]
		if !ok {
			responseMap[serviceKey] = 1
		} else {
			responseMap[serviceKey] = v + 1
		}
	}
}
