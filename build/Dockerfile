# cross-platform build step
ARG GO_VERSION=1.21.1
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build-stage
WORKDIR /usr/src/app
COPY go.mod ./
RUN go mod download && go mod verify
COPY . .
# Compile linux/amd64
FROM --platform=$BUILDPLATFORM build-stage AS build-linux-amd64
RUN CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -o /linux/amd64/server       cmd/server.go
# Compile linux/arm64
FROM --platform=$BUILDPLATFORM build-stage AS build-linux-arm64
RUN CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build -o /linux/arm64/server       cmd/server.go
# Compile windows-amd64
FROM --platform=$BUILDPLATFORM build-stage AS build-windows-amd64
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o /windows/amd64/server.exe cmd/server.go
# package step windows ltsc2019
FROM mcr.microsoft.com/windows/nanoserver:ltsc2019 as package-windows-amd64-ltsc2019
COPY --from=build-windows-amd64 /windows/amd64/server.exe /server.exe
USER ContainerUser
ENTRYPOINT ["/server.exe"]
# package step windows ltsc2022
FROM mcr.microsoft.com/windows/nanoserver:ltsc2022 as package-windows-amd64-ltsc2022
COPY --from=build-windows-amd64 /windows/amd64/server.exe /server.exe
USER ContainerUser
ENTRYPOINT ["/server.exe"]
# package step linux amd64
FROM alpine:latest AS package-linux-amd64-
COPY --from=build-linux-amd64 /linux/amd64/server /server
RUN apk add libcap
RUN setcap CAP_NET_BIND_SERVICE=+eip /server
RUN rm /bin/sh
USER nobody
ENTRYPOINT ["/server"]
# package step linux arm64
FROM alpine:latest AS package-linux-arm64-
COPY --from=build-linux-arm64 /linux/arm64/server /server
RUN apk add libcap
RUN setcap CAP_NET_BIND_SERVICE=+eip /server
RUN rm /bin/sh
USER nobody
ENTRYPOINT ["/server"]
# final step
FROM package-${TARGETOS}-${TARGETARCH}-${WINVER} as FINAL
