FROM golang:1.25.1-alpine3.22 AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o server ./cmd/main.go

FROM alpine:3.22

WORKDIR /app

RUN apk update && apk upgrade

COPY --from=builder /build/server ./server

RUN chmod +x server

EXPOSE 6011

CMD [ "./server" ]