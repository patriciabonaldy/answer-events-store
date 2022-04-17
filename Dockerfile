FROM golang:alpine AS build
RUN apk add --update git
RUN apk add --no-cache bash
RUN apk add build-base
WORKDIR /go/src/github.com/pbonaldy/bequest_challenge
COPY . /go/src/github.com/pbonaldy/bequest_challenge
COPY go.mod /go/src/github.com/pbonaldy/bequest_challenge
COPY go.sum /go/src/github.com/pbonaldy/bequest_challenge

RUN CGO_ENABLED=1 GOOS=linux  go build -ldflags "-linkmode external -extldflags -static"  -o api /go/src/github.com/pbonaldy/bequest_challenge/cmd/main.go
RUN chmod +x api

# Building image with the binary
FROM alpine:latest
EXPOSE 8080
WORKDIR /app
COPY --from=build /go/src/github.com/pbonaldy/bequest_challenge/api /app

ENTRYPOINT ["/app/api"]