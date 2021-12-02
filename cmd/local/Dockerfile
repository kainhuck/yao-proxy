ARG BUILDER_BASE=golang:1.16-alpine
FROM ${BUILDER_BASE} AS builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --update --no-cache make

WORKDIR /yao-proxy

COPY . .

ENV GOPROXY=https://goproxy.cn,direct

RUN make build-local

FROM alpine:3.12

WORKDIR /

COPY --from=builder /yao-proxy/cmd/local/yp-local /bin/
COPY --from=builder /yao-proxy/cmd/local/res/config.json /etc/yao-proxy/config.json

ENTRYPOINT ["/bin/yp-local"]