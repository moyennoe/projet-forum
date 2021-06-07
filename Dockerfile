FROM golang:1.14.2

RUN mkdir /build
WORKDIR /build

RUN  export GO111MODULE=on
RUN go get github.com/gorilla/sessions
RUN go get github.com/mattn/go-sqlite3
RUN cd /build && git clone https://git.ytrack.learn.ynov.com/ACLAVERIA/projet-forum

RUN cd /build/projet-forum && go build

EXPOSE 8080

ENTRYPOINT [ "/build/projet-forum"]