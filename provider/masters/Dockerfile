FROM golang:1.13.0-buster

WORKDIR /go/src

# add apt dependencies 
RUN curl -sL https://deb.nodesource.com/setup_10.x | bash -

# install tools
RUN apt-get update
RUN apt-get -y install zip unzip

# install protocol buffers
RUN curl -OL https://github.com/google/protobuf/releases/download/v3.9.1/protoc-3.9.1-linux-x86_64.zip
RUN unzip protoc-3.9.1-linux-x86_64.zip -d protoc3
RUN mv protoc3/bin/* /usr/local/bin/
RUN mv protoc3/include/* /usr/local/include/
RUN go get -u github.com/golang/protobuf/protoc-gen-go

# install grpc for golang lib
RUN go get -u google.golang.org/grpc

# build synerex daemon
WORKDIR /synerex_simulation
COPY /provider/masters ./provider/masters
COPY /provider/master ./provider/master
COPY /nodeserv ./nodeserv
COPY /server ./server
COPY /provider/agent ./provider/agent
COPY /provider/simutil ./provider/simutil
COPY /util ./util
COPY /api ./api
COPY /nodeapi ./nodeapi

# expose port
#EXPOSE 10000
RUN cd provider/master && go build master-provider.go
RUN cd server && go build synerex-server.go message-store.go
RUN cd nodeserv && go build nodeid-server.go
WORKDIR /synerex_simulation/provider/masters
RUN go build masters-provider.go
ENTRYPOINT [ "./masters-provider" ]