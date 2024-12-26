FROM golang:1.22 AS builder

WORKDIR /app

COPY ../go.* ./

RUN go mod download

COPY ../ ./

RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -v -o favorites ./cmd/favorites

FROM golang:1.22 AS test-runner

WORKDIR /app

COPY ../go.* ./
RUN go mod download

COPY ../ ./

CMD ["go", "test", "-v", "./..."]

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/favorites /favorites

CMD ["/favorites"]