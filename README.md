# Go Service Example
This simple applications uses two services: a Proxy accessed via HTTP and a Processor that serves through gRPC.

The Proxy accepts two integers and calls the Processor (via gRPC) which retruns the multiplication of both numbers.

## Installation

### Clone the repository
```
git clone https://github.com/santiago-tolley/go-service-example.git
```
### Compile Protobuf
```
protoc pb/processor.proto --go_out=plugins=grpc:.
```
Note: install protoc-gen-go with `go get -u github.com/golang/protobuf/protoc-gen-go`

### Create Go module
```
go mod init go-service-example
```

### Build
```
go build cmd/processor/main.go
go build cmd/proxy/main.go
```

## Running the Services
First start the Processor service:
```
go run cmd/processor/main.go
```
Then, once its running, start the Proxy service:
```
go run cmd/proxy/main.go
```
The Processor listens on port 8081 (gRPC), and the Proxy on 8080 (HTTP)


## Calling the Proxy
To call the Proxy run commands as the following:
```
curl --location --request POST 'localhost:8080/multiply/' \
--header 'Content-Type: application/json' \
--data-raw '{ 
	"value": 6, 
	"multiplier": 5 
}'
```
This command targets the "/multiply/" endpoint. 

Set the values of "value" and "multiplier" as the two integers to multiply.

The multiplication in calculated by the Processor service.