FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev openssl

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o sfback .

FROM alpine:latest
RUN apk add --no-cache openssl ca-certificates python3
WORKDIR /app
COPY --from=builder /build/sfback .
COPY --from=builder /build/utilities/c14n.py ./utilities/
RUN mkdir -p /app/data/folders /app/data/keys /app/data/temp

EXPOSE 5001
CMD ["./sfback"]
