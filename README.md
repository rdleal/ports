# Ports Service

This service exposes an HTTP endpoint for receiving a JSON file containing a list of ports
and updates or creates thoses ports to a database.

## Running the service

The entrypoint of the service is located at [cmd/api](./cmd/api) package.

The following command starts the application locally at 8080 HTTP port:

```sh
$ cd cmd/api && PORT=8080 go run main.go
```

## Building Dockerfile

There's a Dockerfile for this service in [build](./build) folder.

The following command builds the application at 8080 HTTP port:

```sh
$ docker build -f build/Dockerfile -t ports-service .
```

The following command starts the container:

```sh
$ docker run -it -p 8080:8080 ports-service
```

## Testing

This service has unit tests for all packages inside [internal](./internal) folder.

### Unit testing

The following command runs all the tests available:

```sh
$ go test -v -cover ./...
```
### Integration testing

The following command runs a manual integration test:

```go
$ curl http://localhost:8080/ports -F "ports=@ports.json" -v
```

where `ports.json` is the path for the file containing a list of ports.
