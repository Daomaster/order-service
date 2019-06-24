# compile the service
FROM golang:1.12 as builder
COPY . /order-service
WORKDIR /order-service
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o order-service

# use scratch for minimal image size
FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /order-service .
ENTRYPOINT ["./order-service"]