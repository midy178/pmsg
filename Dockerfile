FROM golang:1.21.3-alpine3.18 as builder

WORKDIR /build
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s"

FROM --platform=linux/amd64 alpine:latest

# Define the config file name | 定义配置文件名
#ARG CONFIG_FILE=core.yaml
# Define the author | 定义作者
ARG AUTHOR=1228022817@qq.com

LABEL org.opencontainers.image.authors=${AUTHOR}

WORKDIR /app

COPY --from=builder /build/pmsg /usr/bin/

ENTRYPOINT ["pmsg"]
