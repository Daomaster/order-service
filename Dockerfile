FROM golang:1.12 as builder
COPY . /order-service
WORKDIR /order-service
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o order-service

FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /order-service .
ENTRYPOINT ["./order-service"]