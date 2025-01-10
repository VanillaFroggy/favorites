FROM golang:1.23 AS builder

WORKDIR /app

COPY ../go.* ./

RUN go mod download

COPY ../ ./

RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -v -o favorites ./cmd/favorites

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/favorites /favorites

CMD ["/favorites"]