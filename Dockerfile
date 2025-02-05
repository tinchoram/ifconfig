FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./pkg/main.go


FROM alpine:3.19
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/main .

COPY --from=builder /app/views ./views
COPY --from=builder /app/public ./public

EXPOSE 3000

CMD ["./main"]