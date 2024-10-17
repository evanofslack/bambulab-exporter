FROM golang:1.22-alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build cmd/bambu-exporter/main.go

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]

