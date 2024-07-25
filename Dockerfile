FROM golang:1.22.5-alpine AS builder

WORKDIR /dodo

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o dodo
RUN echo "{}" > config.json

FROM gcr.io/distroless/static-debian12:latest

WORKDIR /dodo

COPY --from=builder /dodo/dodo /dodo/dodo
COPY --from=builder /dodo/config.json /dodo/config.json

ENTRYPOINT ["./dodo", "-c", "/dodo/config.json"]
CMD []