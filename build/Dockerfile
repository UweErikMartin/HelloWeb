# package step windows amd64
FROM mcr.microsoft.com/windows/nanoserver:ltsc2019 as package-windows-amd64
COPY windows/amd64/server.exe /server.exe
USER ContainerUser

FROM alpine:latest AS package-linux-amd64
COPY linux/amd64/server /server
RUN apk add libcap
RUN setcap CAP_NET_BIND_SERVICE=+eip /server
RUN chmod 555 /server
RUN ls -al /
# RUN rm /bin/sh
# USER nobody

FROM alpine:latest AS package-linux-arm64
COPY linux/arm64/server /server
RUN apk add libcap
RUN setcap CAP_NET_BIND_SERVICE=+eip /server
RUN chmod 555 /server
RUN ls -al /
# RUN rm /bin/sh
# USER nobody

# final step
FROM package-${TARGETOS}-${TARGETARCH} as FINAL
ENTRYPOINT [ "/server" ]