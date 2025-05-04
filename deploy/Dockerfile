FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o astro .

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/astro .
COPY ephems/de440.bsp ./ephems/de440.bsp

ENTRYPOINT ["./astro"]
