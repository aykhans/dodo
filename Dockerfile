FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o dodo

FROM gcr.io/distroless/static-debian12:latest

WORKDIR /

COPY --from=builder /src/dodo /dodo

ENTRYPOINT ["./dodo"]