FROM ubuntu:16.04
MAINTAINER Dmitry Veselov <d.a.veselov@yandex.ru>

RUN apt-get update
RUN apt-get install -y software-properties-common
RUN add-apt-repository -y ppa:libreoffice/ppa 
RUN apt-get update
RUN apt-get install -y golang git curl
RUN apt-get install -y libreoffice libreofficekit-dev

ENV GOPATH /go

RUN mkdir /go
ADD . /go/src/github.com/docsbox/go-libreofficekit
WORKDIR /go/src/github.com/docsbox/go-libreofficekit
RUN echo "go test -race -coverprofile=coverage.txt -covermode=atomic && bash <(curl -s https://codecov.io/bash) -t 473da5a7-66ec-45dc-b4ed-eb758ce8a66b" > test.sh
