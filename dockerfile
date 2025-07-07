FROM golang:1.18-alpine AS builder

WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o WEBSITE .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/WEBSITE .
COPY frontend/ ./frontend/

EXPOSE 3000
CMD ["./WEBSITE"]