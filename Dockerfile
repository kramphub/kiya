FROM golang:1.9.2
WORKDIR /go/src/github.com/kramphub/kiya/
COPY . .
ARG version
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=$version" .

RUN apk add -U --no-cache ca-certificates

FROM scratch
COPY --from=0 /go/src/github.com/kramphub/kiya .

# see https://github.com/drone/ca-certs/blob/master/Dockerfile
COPY /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/kiya"]