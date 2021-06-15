# Build stage
FROM golang:1.16 AS builder
WORKDIR /go/src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go build -o /go/bin/gochopchop cmd/main.go

# Prod stage
FROM alpine:3.8
COPY --from=builder /go/bin/gochopchop /bin/gochopchop
COPY chopchop.yml /etc/chopchop.yml
ENTRYPOINT [ "/bin/gochopchop" ]
