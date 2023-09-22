# Build the application from source
FROM golang:1.21 AS build-stage

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /helloweb ./...

# Run the tests in the container
# FROM build-stage AS run-test-stage
# RUN go test -v ./...
RUN /helloweb

# Deploy the application binary into a lean image
FROM alpine:3.14 AS build-release-stage

COPY --from=build-stage /helloweb /helloweb

# create a non-root user without home directory
RUN addgroup -g 1000 nonroot
RUN adduser -u 1000 -G nonroot -H -D nonroot
USER nonroot:nonroot

CMD ["/helloweb"]