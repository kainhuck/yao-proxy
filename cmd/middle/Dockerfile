ARG BUILDER_BASE=golang:1.16-alpine
FROM ${BUILDER_BASE} AS builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --update --no-cache make

WORKDIR /yao-proxy

COPY . .

RUN make build-middle

FROM alpine:3.12

WORKDIR /

COPY --from=builder /yao-proxy/cmd/middle/yp-middle /bin/
COPY --from=builder /yao-proxy/cmd/middle/res/config.json /etc/yao-proxy/config.json

ENTRYPOINT ["/bin/yp-middle"]