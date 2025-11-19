FROM golang:1.21-alpine AS builder
WORKDIR /app
# Gerekli paketleri indir
RUN go mod init kobay-app
RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
# Derle
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Küçük imaj oluştur
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]