ARG WINVER
FROM golang:${GO_VERSION} AS build-stage
# package step windows amd64 - LTSC2019
FROM mcr.microsoft.com/windows/nanoserver:${WINVER} as package-windows-amd64
COPY windows/amd64/server.exe /server.exe
USER ContainerUser
ENTRYPOINT [ "/server.exe" ]
# package step linux amd64
FROM alpine:latest AS package-linux-amd64
COPY linux/amd64/server /server
RUN apk add libcap
RUN setcap CAP_NET_BIND_SERVICE=+eip /server
RUN rm /bin/sh
USER nobody
ENTRYPOINT [ "/server" ]
# package step linux arm64
FROM alpine:latest AS package-linux-arm64
COPY linux/arm64/server /server
RUN apk add libcap
RUN setcap CAP_NET_BIND_SERVICE=+eip /server
RUN rm /bin/sh
USER nobody
ENTRYPOINT [ "/server" ]
# final step
FROM package-${TARGETOS}-${TARGETARCH} as FINAL
