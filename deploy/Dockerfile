FROM golang:1.22-bookworm AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/app ./cmd/$TARGET

FROM gcr.io/distroless/static
COPY --from=builder /bin/app /app
ENTRYPOINT ["/app"]
