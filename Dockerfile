FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build
WORKDIR /src
COPY . .
ARG VERSION
ARG TARGETOS
ARG TARGETARCH
RUN mkdir -p /out && \
    cd /src/cmd/kiya && \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build \
      -a \
      -ldflags "-s -w -X 'main.version=$VERSION'" \
      -o /out/kiya

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=build /out/kiya /usr/bin/kiya
RUN chmod 755 /usr/bin/kiya
ENTRYPOINT ["/usr/bin/kiya"]
