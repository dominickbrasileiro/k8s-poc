FROM golang:1.19.1-alpine as builder

WORKDIR /app/disperser

COPY ./disperser .

RUN go build -o main


FROM alpine:3.15

WORKDIR /app/disperser

COPY --from=builder /app/disperser/main .

CMD ["./main"]