FROM golang:1.22 AS builder
WORKDIR /source
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o se-vende ./cmd/bot/main.go

FROM alpine:latest
RUN apk --no-cache add tzdata ca-certificates
WORKDIR /bot/
COPY --from=builder /source/se-vende .
CMD ["./se-vende"]
