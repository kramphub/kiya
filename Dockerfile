FROM golang:1.9.2
WORKDIR /go/src/github.com/kramphub/kiya/
COPY . .
ARG version
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=$version" .

FROM scratch
COPY --from=0 /go/src/github.com/kramphub/kiya .
ENTRYPOINT ["/kiya"]