FROM golang:1.10.5-stretch

ENV LIBRDKAFKA_VERSION 0.11.5

ADD . /go/src/git.ng.bluemix.net/ageiger/kafka-test

RUN apt-get -y update \
  && apt-get install -y upx-ucl zip libssl-dev libsasl2-dev ca-certificates wget jq zip

RUN curl -Lk -o /root/librdkafka-${LIBRDKAFKA_VERSION}.tar.gz https://github.com/edenhill/librdkafka/archive/v${LIBRDKAFKA_VERSION}.tar.gz && \
  tar -xzf /root/librdkafka-${LIBRDKAFKA_VERSION}.tar.gz -C /root && \
  cd /root/librdkafka-${LIBRDKAFKA_VERSION} && \
  ./configure --prefix /usr && make && make install && make clean && ./configure --clean

RUN apt-get update \
  && cd /go/src/git.ng.bluemix.net/ageiger/kafka-test \
  && CGO_ENABLED=1 make build/kafka-test \
  && ln -s /go/src/git.ng.bluemix.net/ageiger/kafka-test/build/kafka-test /kafka-test \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

ENTRYPOINT [ "/kafka-test" ]
