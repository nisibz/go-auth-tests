FROM golang:1.24-bookworm

WORKDIR /app

RUN go install github.com/air-verse/air@latest

CMD [ "air" ]
