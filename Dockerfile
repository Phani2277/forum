# Build stage
FROM golang:1.21 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o forum ./cmd/web

# Run stage
FROM debian:stable-slim
WORKDIR /app
COPY --from=build /app/forum .
COPY templates ./templates
COPY static ./static
EXPOSE 8080
CMD ["./forum"]
