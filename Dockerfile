FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .
RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

# Create a directory for the service account key
RUN mkdir -p /etc/secrets

# The actual key file should be mounted at runtime
ENV GOOGLE_APPLICATION_CREDENTIALS=/etc/secrets/service-account-key.json

EXPOSE 8070

CMD ["./main"]
