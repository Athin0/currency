FROM golang:latest

WORKDIR /currency
COPY go.* ./
RUN go mod download

COPY ./ /currency

RUN apt-get update && apt-get -y upgrade

RUN go build -o myapp ./cmd/main.go

CMD ["/currency/myapp"]