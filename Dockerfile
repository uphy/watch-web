FROM golang:1.13 as server-builder
WORKDIR /go/src/github.com/uphy/watch-web/server
COPY server .
RUN go get -u github.com/markbates/pkger/cmd/pkger && pkger -o resources
RUN CGO_ENABLED=0 GOOS=linux go build -o /server
RUN CGO_ENABLED=0 GOOS=linux go get -u github.com/itchyny/gojq/cmd/gojq

FROM node:8.11 as client-builder
WORKDIR /app
COPY client .
RUN yarn
RUN yarn build

FROM alpine
RUN apk add --no-cache ca-certificates curl && update-ca-certificates
WORKDIR /app
COPY --from=server-builder /server .
COPY --from=server-builder /go/bin/gojq /usr/bin/
COPY server/config.yml .
COPY --from=client-builder /app/dist ./static
COPY server/scripts/ ./scripts/
EXPOSE 8080
ENTRYPOINT ["./server"]