# SimpleCountingService
A simple service responding with request counts to explore invocation patterns, or test concurrent use of multiple service instances (e.g., to explore load balancing). To support the latter case, the service automatically creates a unique identifier (and background color for the HTML output) to distinguish instances at runtime. For diagnostic purposes the service optionally prints invocations to the console.

By default, the service listens on port 8080.

Endpoints include:

* `:8080/count` to continue counting invocations. Note that handler will sleep to one second before returning (to simulate delay due to processing workload).
* `:8080/reset` to reset the count
* `:8080/exit` to stop the service without OS error (status code 0)
* `:8080/kill` to stop the service with OS error (status code 1)

This service comes with a Dockerfile for deployment. To build, call

`docker build -t countingservice .`

* Note: if docker is not running with sudo privileges, prefix the command with `sudo`

For deployment, call

`docker run -d -p 8080:8080 --name runningService --restart=always countingservice`

Alternatively, the service can be started using docker compose (`docker compose up -d`). Please see the parameterization in `compose.yml`.

The repository further includes a complementary client that can be used to evaluate invocation patterns (e.g., errors, distribution across instances in the case of load balancing).
