FROM node:10 as frontend-builder
WORKDIR /app
COPY frontend .
RUN yarn
RUN yarn build

FROM golang:1.13-alpine as server-builder
WORKDIR /go/src/github.com/uphy/watch-web
# Build app
COPY . .
COPY --from=frontend-builder /app/dist/ ./frontend/dist/
RUN go get -u github.com/markbates/pkger/cmd/pkger && pkger -o resources
RUN CGO_ENABLED=0 GOOS=linux go build -o /server
# Build gojq
RUN CGO_ENABLED=0 GOOS=linux go get -u github.com/itchyny/gojq/cmd/gojq

FROM alpine
RUN apk add --no-cache ca-certificates curl && update-ca-certificates
WORKDIR /app
COPY --from=server-builder /server .
COPY --from=server-builder /go/bin/gojq /usr/bin/
COPY config.yml .
COPY scripts/ ./scripts/
EXPOSE 8080
ENTRYPOINT ["./server"]