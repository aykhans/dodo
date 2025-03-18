FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags "-s -w" -o dodo
RUN echo "{}" > config.json

FROM gcr.io/distroless/static-debian12:latest

WORKDIR /

COPY --from=builder /src/dodo /dodo
COPY --from=builder /src/config.json /config.json

ENTRYPOINT ["./dodo", "-f", "/config.json"]