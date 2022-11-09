FROM golang:1.19.3 as builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /go/src/github.com/wadeling/kafka-demo
COPY . ./
RUN go build -o kafka_bench .

FROM ubuntu:22.04 
COPY --from=builder /go/src/github.com/wadeling/kafka-demo/kafka_bench /usr/local/bin/kafka_bench
COPY ./entrypoint.sh /usr/local/bin/entrypoint.sh
ENTRYPOINT ["entrypoint.sh"]

