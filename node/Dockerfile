FROM golang:1.19.1-alpine as builder

WORKDIR /app/node

COPY ./node .

RUN go build -o main


FROM alpine:3.15

WORKDIR /app/node

COPY --from=builder /app/node/main .

EXPOSE 32000

CMD ["./main"]