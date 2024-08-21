FROM golang:1.18-alpine3.16 as builder
WORKDIR /go/src/k8s-web/server
COPY . .

RUN go env -w GO111MODULE=on \
   && go env -w GOPROXY=https://goproxy.cn,direct \
   && go env -w CGO_ENABLED=0 \
   && go env \
   && go mod tidy \
   && go build -o server .

FROM alpine:latest

LABEL MAINTAINER="muxian@imooc.com"

WORKDIR /go/src/k8s-web/server
COPY --from=0 /go/src/k8s-web/server/config.yaml ./config.yaml
COPY --from=0 /go/src/k8s-web/server/.kube/config ./.kube/config
COPY --from=0 /go/src/k8s-web/server/server ./
EXPOSE 8082
ENTRYPOINT ./server