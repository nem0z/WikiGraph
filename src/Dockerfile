FROM golang:latest

WORKDIR /app

COPY ./src ./src

COPY .env .env

WORKDIR /app/src

RUN go mod download

RUN go build -o wikigraph .

CMD ["./wikigraph"]