FROM ubuntu
MAINTAINER LIZHIXIANG
WORKDIR /var/cache/apt/archives
RUN apt-get update
RUN apt-get install python --assume-yes
RUN apt-get install apt-utils --assume-yes
RUN apt-get install libaio1 --assume-yes
RUN apt-get install libibverbs1 --assume-yes
RUN apt-get install librados2 --assume-yes
RUN apt-get install librbd1 --assume-yes
RUN apt-get install librdmacm1 --assume-yes
COPY fio_2.2.10-1ubuntu1_amd64.deb .
RUN dpkg -i fio_2.2.10-1ubuntu1_amd64.deb
RUN apt-get update
WORKDIR /
RUN mkdir -p go/src/fioProject
WORKDIR /go/src/fioProject
RUN mkdir fioTool
RUN mkdir report
COPY fioProject.go .
RUN mv fioProject.go main.go
COPY fioProject .
WORKDIR /go/src/fioProject/fioTool
COPY fioTool.go .
RUN mv fioTool.go main.go
COPY fioTool .
WORKDIR /go/src/fioProject
CMD ./fioProject

