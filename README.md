# Go Stress Test

Simple Stress Test tool for HTTP requests.

## Options

Available options are:
- `url` (**required**): URL to send requests;
- `requests`: Number of requests to send;
- `concurrency`: Number of concurrent requests;
- `method`: HTTP method to use;

## Usage

### Directly run

```bash
# Execute the application
go run cmd/stress_test/main.go -url=<url> -requests=<requests> -concurrency=<workers>
```

### Docker run

```bash
# Build the image
docker build -t go-stress-test .

# Execute the application
docker run go-stress-test --url=<url> --requests=<requests> --concurrency=<workers>
```

## Example API

For testing purposes an Example API was developed at `cmd/mock_server`.

This server will return randomly one of the following status codes with a
random delay between 0 and 2 seconds:
- 500: Internal Server Error;
- 503: Service Unavailable;
- 404: Not Found;
- 401: Unauthorized;
- 200: OK;

To run the server execute:
```bash
go run cmd/mock_server/main.go
```

To run the stress tests against the server execute:
```bash
go run cmd/stress_test/main.go -url=http://localhost:8080 -requests=100 -concurrency=10
```

Example output of the stress test against the Example API:
```bash
$ go run cmd/stress_test/main.go --url=http://localhost:8080 --requests=100 --concurrency=10
Configuration:
- Concurrency: 10
- Requests: 100
- Timeout: 5s
- Method: GET
- URL: http://localhost:8080

Results:
- 500: 21/100 (21%)
- 503: 19/100 (19%)
- 401: 20/100 (20%)
- 200: 18/100 (18%)
- 404: 22/100 (22%)
Total tasks executed: 100

Total time: 10.03s
```

## Testing

To execute all the unit tests run `go test ./...`.
