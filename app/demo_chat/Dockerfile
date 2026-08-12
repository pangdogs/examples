FROM golang:1.25-alpine AS builder

ARG GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0 \
    GOPROXY=${GOPROXY}

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY app/demo_chat/consts ./app/demo_chat/consts
COPY app/demo_chat/server ./app/demo_chat/server

RUN go build -trimpath -ldflags="-s -w" -o /out/chat-server ./app/demo_chat/server

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /out/chat-server ./chat-server
COPY app/demo_chat/bin/cli.pub app/demo_chat/bin/serv.pem ./

EXPOSE 8080 9090

ENTRYPOINT ["./chat-server"]
CMD ["--cli_pub_key=./cli.pub", "--serv_priv_key=./serv.pem"]
