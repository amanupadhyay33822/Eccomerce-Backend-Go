# Stage 1: Build the Go binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

# copy go mod files
COPY go.mod go.sum ./

RUN go mod download

# copy source code
COPY . .

# build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Small production image
FROM alpine:latest

WORKDIR /root/

# copy binary from builder
COPY --from=builder /app/main .



# expose port
EXPOSE 8080

# run application
CMD ["./main"]