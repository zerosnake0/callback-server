# Build stage
FROM golang:1.13.7-alpine3.11 AS build

ENV GOPATH=/go \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk update \
 && apk add ca-certificates

WORKDIR $GOPATH/src/callback-server

COPY vendor vendor

COPY main.go main.go

COPY pkg pkg

RUN go build -a -installsuffix cgo -ldflags="-s -w" \
    -o /bin/server main.go

# Production stage
FROM golang:1.13.7-alpine3.11

ENV LANG=C.UTF-8

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /bin

COPY --from=build /bin/server .

EXPOSE 80

ENTRYPOINT ["/bin/server"]
