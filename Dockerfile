FROM golang:1.22.1 AS builder
WORKDIR /filmoteka
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cmd/filmoteka/bin/main ./cmd/filmoteka/
FROM alpine:latest
WORKDIR /filmoteka
COPY --from=builder /filmoteka/cmd/filmoteka/bin/main filmoteka
CMD ["/filmoteka/filmoteka"]