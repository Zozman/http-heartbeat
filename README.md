# HTTP Heartbeat

A simple container to send a HTTP GET request to a specified url on an interval with the ability to send only if a 200 is returned from another URL.  Useful as a [Kubernetes Sidecar](https://kubernetes.io/docs/concepts/workloads/pods/sidecar-containers/) for reporting health information to an external service.

## Running Locally

You can either run this as a containerized application using [Docker](https://www.docker.com/) or as a direct Go app.  We recommend using the Docker method for ease of use if you have [Docker Desktop](https://www.docker.com/products/docker-desktop/) installed.

### Docker

1. Make sure [Docker Desktop](https://www.docker.com/products/docker-desktop/) is installed
2. Set your environmental variables inside the [docker-compose.yaml](docker-compose.yaml) file
3. Run `docker compose up --build`

### Go

1. Run `go mod download` in the `src` directory to install dependencies
2. Run `go build -o ./http_heartbeat ./cmd/main.go` in the `src` directory to build the binary
3. Run `./http_heartbeat` inm the `src` directory while passing settings as environmental variables

## Environmental Variables

- `HEARTBEAT_URL`
    - **Required**
    - URL to send GET request to on interval
- `TEST_URL`
    - If set, app will first perform a HTTP GET request to this URL and require an OK status before sending the heartbeat to the `HEARTBEAT_URL`
- `INTERVAL`
    - Default Value: `30`
    - Number of seconds between heartbeat cycles
- `VERBOSE`
    - If set, enables verbose logging