# Build the application from source
FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod ./
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /helloweb

# Run the tests in the container
# FROM build-stage AS run-test-stage
# RUN go test -v ./...

# Deploy the application binary into a lean image
FROM alpine:3.14 AS build-release-stage

WORKDIR /

COPY --from=build-stage /helloweb /helloweb

EXPOSE 8080

# create a non-root user without home directory
RUN addgroup -g 1000 nonroot
RUN adduser -u 1000 -G nonroot -H -D nonroot
USER nonroot:nonroot

ENTRYPOINT ["/helloweb"]