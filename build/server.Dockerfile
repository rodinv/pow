FROM golang:1.17 AS builder

WORKDIR /app
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o server cmd/server/main.go

FROM alpine:3.12.0 AS launcher

WORKDIR /
COPY --from=builder /app/server .

CMD /server
