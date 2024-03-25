FROM golang:1.22 AS builder
WORKDIR /source
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o se-vende .

FROM scratch
COPY --from=builder /source/se-vende .
CMD ["./se-vende"]
