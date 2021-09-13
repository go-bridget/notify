FROM golang:1.16

# install protobuf
ENV PB_VER 3.10.1
ENV PB_URL https://github.com/google/protobuf/releases/download/v${PB_VER}/protoc-${PB_VER}-linux-x86_64.zip

RUN apt-get -qq update && apt-get -qqy install curl git make unzip gettext rsync

RUN mkdir -p /tmp/protoc && \
    curl -L ${PB_URL} > /tmp/protoc/protoc.zip && \
    cd /tmp/protoc && \
    unzip protoc.zip && \
    cp /tmp/protoc/bin/protoc /usr/local/bin && \
    cp -R /tmp/protoc/include/* /usr/local/include && \
    chmod go+rx /usr/local/bin/protoc && \
    cd /tmp && \
    rm -r /tmp/protoc

# Get the source from GitHub
RUN go get -u google.golang.org/grpc

# Install protoc-gen-go
RUN go get -u github.com/golang/protobuf/protoc-gen-go

# Install protoc-gen-twirp
RUN go get -u github.com/twitchtv/twirp/protoc-gen-twirp

# Install googleapis proto files, grpc-gateway generator (http to grpc mapping)

RUN git clone https://github.com/googleapis/googleapis /opt/googleapis
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

# Install mig to handle migrations
RUN git clone https://github.com/go-bridget/mig /opt/mig
RUN cd /opt/mig && make build && cp _build/mig /usr/local/bin

RUN go get -u github.com/fatih/faillint

RUN go get -u github.com/go-bridget/twirp-swagger-gen
