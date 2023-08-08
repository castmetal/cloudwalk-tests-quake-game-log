FROM golang:1.20.4 AS builder
WORKDIR /src/reader_log_script
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p ./build/bin
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/bin/ ./cmd/...

FROM alpine:latest AS script
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /src/reader_log_script/build/bin/ .
CMD ["./reader_log_script"]