FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o todo src/main.go

FROM alpine:3.16
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
