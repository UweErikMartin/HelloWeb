ARG GO_VERSION
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build-stage
WORKDIR /usr/src/app
# copy the dependency files and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -o /linux/amd64/server       cmd/server.go
RUN CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build -o /linux/arm64/server       cmd/server.go
RUN CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o /windows/amd64/server.exe cmd/server.go
# package step windows amd64 - LTSC2019
FROM mcr.microsoft.com/windows/nanoserver:ltsc2019 as package-windows-amd64
ARG TARGETOS TARGETARCH TARGETVARIANT
RUN echo $TARGETVARIANT
COPY --from=build-stage /windows/amd64/server.exe /server.exe
USER ContainerUser
ENTRYPOINT [ "/server.exe" ]
# package step linux amd64
FROM alpine:latest AS package-linux-amd64
ARG TARGETOS TARGETARCH TARGETVARIANT
RUN echo $TARGETVARIANT
COPY --from=build-stage /linux/amd64/server /server
RUN apk add libcap
RUN setcap CAP_NET_BIND_SERVICE=+eip /server
RUN rm /bin/sh
USER nobody
ENTRYPOINT [ "/server" ]
# package step linux arm64
FROM alpine:latest AS package-linux-arm64
ARG TARGETOS TARGETARCH TARGETVARIANT
RUN echo $TARGETVARIANT
COPY --from=build-stage /linux/arm64/server /server
RUN apk add libcap
RUN setcap CAP_NET_BIND_SERVICE=+eip /server
RUN rm /bin/sh
USER nobody
ENTRYPOINT [ "/server" ]
# final step
FROM package-${TARGETOS}-${TARGETARCH} as FINAL
