FROM golang
MAINTAINER LIZHIXIANG
RUN apt-get update
RUN apt-get install fio --assume-yes
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

