FROM golang:1.19.1 AS builder

RUN mkdir /app
COPY app/ /app/
WORKDIR /app
RUN CGO_ENABLED=0 go build -o main main.go
RUN wget https://tuf.werf.io/targets/releases/1.2.240/linux-arm64/bin/werf -O /usr/local/bin/werf
RUN chmod a+x /usr/local/bin/werf

FROM alpine:3.16.2
RUN apk add curl
RUN mkdir /app
COPY --from=builder /app /app
COPY --from=builder /usr/local/bin/werf /usr/local/bin/werf
WORKDIR /app
ENTRYPOINT ["/app/main"]
