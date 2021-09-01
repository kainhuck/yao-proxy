ARG BUILDER_BASE=golang:1.16-alpine
FROM ${BUILDER_BASE} AS builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN  --mount=type=cache,target=/var/cache/apk apk add --update --no-cache make

WORKDIR /yao-proxy

COPY . .

RUN make build-local

FROM alpine:3.12

EXPOSE 20808

WORKDIR /

COPY --from=builder /yao-proxy/cmd/local/yp-local /bin/
COPY --from=builder /yao-proxy/cmd/local/res/config.json /etc/yao-proxy/config.json

ENTRYPOINT ["/bin/yp-local"]