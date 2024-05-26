FROM golang:latest

WORKDIR /client

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 80

CMD ["./main"]