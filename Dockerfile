FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o bin main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/ .
EXPOSE 8000
CMD ["/app/bin"]
