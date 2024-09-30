FROM golang:1.22-alpine

WORKDIR /app

COPY main.go /app/

RUN go build -v -o server main.go

EXPOSE 3333

CMD [ "/app/server" ]