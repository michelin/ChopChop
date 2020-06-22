FROM golang:1.13 AS build
RUN mkdir /app
ADD . /app/
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY chopchop.yml ./
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build .
CMD ["/app/gochopchop"]

FROM alpine:3.8
RUN mkdir -p /tmp
COPY --from=build /app/gochopchop /tmp/gochopchop
COPY --from=build /app/chopchop.yml /tmp/chopchop.yml
CMD ["/tmp/gochopchop"]
