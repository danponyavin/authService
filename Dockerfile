FROM golang:1.22-alpine as builder

COPY . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix 'static' -o /app cmd/main.go

FROM alpine:latest
COPY --from=builder /app .

EXPOSE 8000/tcp
ENTRYPOINT ["/app"]