FROM golang:latest

ENV GOPROXY = http://goproxy.cn

COPY . /app

WORKDIR /app

RUN go build -o server .

EXPOSE 1080

CMD [ "server" ]