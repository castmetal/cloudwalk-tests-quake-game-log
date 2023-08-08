FROM golang:1.20.4 AS builder
WORKDIR /src/reader_log_script
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p ./build/bin
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/bin/ ./cmd/...

FROM alpine:latest AS script
WORKDIR /
COPY --from=builder /src/reader_log_script/build/bin/ .
ENTRYPOINT ["./reader_log_script reader_log_script --execute=true"]