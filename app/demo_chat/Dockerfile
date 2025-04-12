FROM golang:1.23-alpine

ENV GOPROXY=https://goproxy.cn,direct
ENV ETCD_ADDR=${ETCD_ADDR}
ENV NATS_ADDR=${NATS_ADDR}

WORKDIR /app

COPY ./bin ./bin
COPY ./cli ./cli
COPY ./server ./server

WORKDIR /app/server

RUN go mod init server
RUN go mod tidy
RUN go build -o ../bin/server

WORKDIR /app/cli

RUN go mod init cli
RUN go mod tidy
RUN go build -o ../bin/cli

WORKDIR /app/bin

RUN echo -e '#!/bin/sh\n\
exec ./server --etcd.address=${ETCD_ADDR} --nats.address=${NATS_ADDR}\n\
'> ./entrypoint.sh
RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
