FROM golang:1.23.3-alpine AS builder

WORKDIR /app

COPY . ./

RUN go build -o frogfoot-server ./cmd/server/main.go

FROM golang:1.23.3-alpine

RUN apk add sqlite

COPY --from=builder /app/frogfoot-server /app/frogfoot-server

CMD ["/app/frogfoot-server"]
