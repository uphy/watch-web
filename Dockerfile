FROM golang:1.10 as server-builder
WORKDIR /go/src/github.com/uphy/watch-web/server
COPY server .
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -o /server

FROM node:8.11 as client-builder
WORKDIR /app
COPY client .
RUN yarn
RUN yarn build

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR /app
COPY --from=server-builder /server .
COPY server/config.yml .
COPY --from=client-builder /app/dist ./static
EXPOSE 8080
ENTRYPOINT ["./server"]