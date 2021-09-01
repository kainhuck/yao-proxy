ARG BUILDER_BASE=golang:1.16-alpine
FROM ${BUILDER_BASE} AS builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN  --mount=type=cache,target=/var/cache/apk apk add --update --no-cache make

WORKDIR /yao-proxy

COPY . .

RUN make build-remote

FROM alpine:3.12

WORKDIR /

COPY --from=builder /yao-proxy/cmd/remote/yp-remote /bin/
COPY --from=builder /yao-proxy/cmd/remote/res/config.json /etc/yao-proxy/config.json

ENTRYPOINT ["/bin/yp-remote"]