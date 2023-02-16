FROM golang:latest

WORKDIR /go/src/app

# Build and install Service 1
COPY server-1/ .
RUN go get -d -v ./...
RUN go install -v ./...

# Build and install Service 2
COPY server-2/ .
RUN go get -d -v ./...
RUN go install -v ./...

# Start both services
CMD ["sh", "-c", "server-1 & server-2"]
