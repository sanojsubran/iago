FROM golang:alpine

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN chmod 777 go.mod go.sum

RUN ls -la

RUN go mod download

RUN go build -o hnservice .

CMD ["/app/hnservice"]