FROM golang:1.23-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0, GOOS=linux GOARCH=amd64
RUN go build -o test cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/test .
COPY ./migrations ./migrations
RUN chmod +x wait-for-it.sh
CMD ["./test", "--debug"]
EXPOSE 8081