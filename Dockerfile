FROM golang:1.19.4-bullseye as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o bot .

FROM alpine:3.17
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/bot ./
ENTRYPOINT ["./bot"]

