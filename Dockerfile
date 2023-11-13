FROM golang:latest

COPY ./ /app

RUN export GOPATH=/app

WORKDIR /app

RUN go mod download

RUN go build -o main .

ENTRYPOINT [ "./main" ]