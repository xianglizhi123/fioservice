FROM golang
MAINTAINER LIZHIXIANG
WORKDIR /var/cache/apt/archives
RUN apt-get update
RUN apt-get install apt-utils --assume-yes
RUN apt-get install libaio1 --assume-yes
RUN apt-get install libibverbs1 --assume-yes
RUN apt-get install librados2 --assume-yes
RUN apt-get install librbd1 --assume-yes
RUN apt-get install librdmacm1 --assume-yes
COPY fio_2.2.10-1ubuntu1_amd64.deb .
RUN dpkg -i fio_2.2.10-1ubuntu1_amd64.deb
RUN apt-get update
WORKDIR $GOPATH
RUN go get "github.com/gorilla/mux"
WORKDIR $GOPATH/src/fioProject
RUN mkdir fioTool
RUN mkdir report
COPY fioserver.go .
RUN mv fioserver.go main.go
RUN go build
WORKDIR $GOPATH/src/fioProject/fioTool
COPY fioagent.go .
RUN mv fioagent.go main.go
RUN go build
WORKDIR $GOPATH/src/fioProject
CMD ./fioProject

