FROM golang:alpine

WORKDIR /app

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o apiserver main.go

FROM alpine:latest

COPY --from=0 /app ./

EXPOSE 8080

CMD ["./apiserver"]