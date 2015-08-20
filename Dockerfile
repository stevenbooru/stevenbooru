FROM golang:1.4.2

ENV GBUPDATED 08-19-2015
RUN go get github.com/constabulary/gb/...

ADD . /app
WORKDIR /app
RUN mkdir data
RUN gb build all

EXPOSE 6606
