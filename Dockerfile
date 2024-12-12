FROM golang:1.23.4-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o todo src/main.go

FROM alpine:3.21
RUN apk add --no-cache tzdata
WORKDIR /todo/app
COPY --from=builder /app/todo .
WORKDIR /todo/templates
COPY --from=builder /app/templates .
WORKDIR /todo
COPY --from=builder /app/app.env .

EXPOSE 8080
WORKDIR /todo/app
CMD ["./todo"]
