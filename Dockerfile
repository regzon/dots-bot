FROM golang:latest

WORKDIR /src

COPY . .

RUN go mod download

RUN go build -o /usr/bin/dots-game ./cmd/dots-game

CMD ["dots-game"]
