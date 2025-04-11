FROM golang:1.23-alpine

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY ./bin ./bin
COPY ./server ./src

RUN echo -e '#!/bin/sh\n\
set -e\n\
cd /app/src\n\
go mod init chat\n\
go mod tidy\n\
go build -o ../bin/server\n\
cd /app/bin\n\
./server --etcd.address=${ETCD_ADDR} --nats.address=${NATS_ADDR} &\n\
wait -n\n\
exit $?\n\
'> ./bin/start.sh

RUN chmod +x ./bin/start.sh
RUN dos2unix ./bin/start.sh

ENTRYPOINT ["/app/bin/start.sh"]
