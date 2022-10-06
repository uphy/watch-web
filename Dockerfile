# Dockerfile for Heroku deployment

FROM golang:1.19 as server-builder
WORKDIR /go/src/github.com/uphy/watch-web
# Build app
COPY . .
RUN go get -u github.com/markbates/pkger/cmd/pkger && pkger -o pkg/resources
RUN CGO_ENABLED=0 GOOS=linux go build -o /watch-web
# Build gojq
RUN CGO_ENABLED=0 GOOS=linux go get github.com/itchyny/gojq/cmd/gojq

FROM ubuntu:18.04
RUN apt-get update -y && \
    apt-get install -y curl gzip
WORKDIR /app
COPY --from=server-builder /watch-web .
COPY --from=server-builder /go/bin/gojq /usr/bin/
COPY config.yml .
COPY scripts/ ./scripts/
COPY includes/ ./includes/
EXPOSE 8080
CMD ["./watch-web", "start", "--api", "--no-schedule"]
