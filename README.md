# SimpleCountingService
A simple service responding with request counts to explore invocation patterns, or test concurrent use of multiple service instances (e.g., to explore load balancing). To support the latter case, the service automatically creates a unique identifier (and background color for the HTML output) to distinguish instances at runtime.

By default, the service listens on port 8080.

Endpoints include:

* `:8080/count` to continue counting invocations
* `:8080/reset` to reset the count
* `:8080/exit` to stop the service without OS error (status code 0)
* `:8080/kill` to stop the service with OS error (status code 1)

This service comes with a Dockerfile for deployment. To build, call

`docker build -t countingservice .`

For deployment, call

`docker run -d -p 8080:8080 --name runningService --restart=always countingservice`

The repository further includes a complementary client that can be used to evaluate invocation patterns (e.g., errors, distribution across instances in the case of load balancing).
