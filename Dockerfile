# Dockerfile for Heroku deployment

FROM golang:1.19 as server-builder
WORKDIR /go/src/github.com/uphy/watch-web
# Build app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /watch-web

FROM ubuntu:18.04
RUN apt-get update -y && \
    apt-get install -y curl gzip
WORKDIR /app
COPY --from=server-builder /watch-web .
COPY config.yml .
COPY scripts/ ./scripts/
COPY includes/ ./includes/
EXPOSE 8080
CMD ["./watch-web", "start", "--api", "--no-schedule"]
