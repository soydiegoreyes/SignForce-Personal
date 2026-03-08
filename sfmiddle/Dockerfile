FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev openssl

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o sfmiddle .

FROM alpine:latest
RUN apk add --no-cache openssl ca-certificates python3
WORKDIR /app
COPY --from=builder /build/sfmiddle .
COPY --from=builder /build/templates/ ./templates/
COPY --from=builder /build/utilities/insert_qr.py ./utilities/
RUN mkdir -p /app/data/folders /app/data/keys /app/data/temp /app/data/docs /app/data/images /app/data/templates /app/data/emails

EXPOSE 5002
CMD ["./sfmiddle"]
