# Dockerfile for Heroku deployment

FROM node:10 as frontend-builder
WORKDIR /app
COPY frontend .
RUN yarn
RUN yarn build

FROM golang:1.13-stretch as server-builder
WORKDIR /go/src/github.com/uphy/watch-web
# Build app
COPY . .
COPY --from=frontend-builder /app/dist/ ./frontend/dist/
RUN go get -u github.com/markbates/pkger/cmd/pkger && pkger -o pkg/resources
RUN CGO_ENABLED=0 GOOS=linux go build -o /watch-web
# Build gojq
RUN CGO_ENABLED=0 GOOS=linux go get -u github.com/itchyny/gojq/cmd/gojq

FROM ubuntu:18.04
RUN apt-get update -y && \
    apt-get install -y curl gzip
WORKDIR /app
COPY --from=server-builder /watch-web .
COPY --from=server-builder /go/bin/gojq /usr/bin/
COPY config.yml .
COPY scripts/ ./scripts/
EXPOSE 8080
CMD ["./watch-web", "start", "--api", "--no-schedule"]
