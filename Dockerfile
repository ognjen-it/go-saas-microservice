FROM golang:1.19.1 as builder

WORKDIR /go-microservice/

COPY . .

RUN CGO_ENABLED=0 go build -o microservice /go-microservice/main.go

FROM alpine:latest

WORKDIR /go-microservice

COPY --from=builder /go-microservice/ /go-microservice/

EXPOSE 8080

CMD ./microservice