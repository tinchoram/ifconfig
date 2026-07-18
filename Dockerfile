# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath -ldflags="-s -w" \
    -o /app/ifconfig ./cmd/ifconfig

FROM alpine:3.21

# Run as an unprivileged user; the binary never needs root.
RUN adduser -D -H -u 10001 appuser

WORKDIR /app
COPY --from=builder /app/ifconfig ./ifconfig
COPY --from=builder /app/views ./views
COPY --from=builder /app/public ./public

USER appuser

EXPOSE 3000

# busybox wget (bundled with alpine) probes the liveness endpoint.
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://127.0.0.1:3000/status || exit 1

ENTRYPOINT ["./ifconfig"]
