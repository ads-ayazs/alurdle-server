FROM golang:1.16-alpine

WORKDIR /app

COPY ./app ./

RUN go mod download

RUN go build -o /wordle-master

EXPOSE 8080

CMD [ "/wordle-master" ]
