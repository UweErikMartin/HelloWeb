# Get the golang version passed from the caller as it is derived
# out of the go,mod file in the repository. The default is 
# golang version 1.21
ARG GO_VERSION=1.21

# Build the application from source.
# use the official golang docker image
# to compile the source code
FROM golang:${GO_VERSION} AS build-stage

# output the target platform
RUN echo ${TARGETPLATFORM}
# How to use?
# see description in https://hub.docker.com/_/golang
WORKDIR /usr/src/app

# copy the dependency files and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# copy the rest of the source code
# and compile for linux. Put the result
# into the /usr/local/bin standard directory
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/server cmd/server.go

# Run the tests in the container
# FROM build-stage AS run-test-stage
# RUN go test -v ./...

# Deploy the application binary into a lean image
# based on alpine. We don't want to ship the golang
# environment
FROM alpine:latest AS build-release-stage

# copy just the binary into the deployment image
COPY --from=build-stage /usr/local/bin/server /server

# Allow network access to the application. Without 
# it, the container will not run in the strict environment
# of kubernetes
RUN apk add libcap
RUN setcap CAP_NET_BIND_SERVICE=+eip /server

# create a non-root user without home directory. We
# don't want to run our container with root privileges
RUN addgroup -g 1000 nonroot
RUN adduser -u 1000 -G nonroot -H -D nonroot

# run the container with the new user
USER nonroot:nonroot

# Entrypoint
ENTRYPOINT ["/server"]