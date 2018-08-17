FROM golang:1.10 as builder

WORKDIR /go/src/github.com/uphy/watch-web
COPY . .
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -o /watch-web

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR /app
COPY --from=builder /watch-web .
COPY config.yml .
EXPOSE 8080
ENTRYPOINT ["./watch-web"]