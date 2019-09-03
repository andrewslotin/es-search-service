FROM golang:1.12-alpine AS build-env

RUN apk update && apk add --no-cache git
WORKDIR /build

COPY ./ ./
RUN go build -o server -mod vendor

FROM alpine

RUN adduser -D -g '' service
USER service

COPY --from=build-env /build/server /app/server

ENTRYPOINT ["/app/server"]
