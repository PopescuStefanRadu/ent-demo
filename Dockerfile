FROM golang:1.21.5-alpine3.18 as builder

RUN apk add --no-cache git make build-base

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY Makefile ./
COPY pkg ./pkg
COPY cmd ./cmd


RUN make build

FROM alpine:3.18

EXPOSE 8080/tcp

COPY --from=builder /build/server /server

ENTRYPOINT ["/server"]
