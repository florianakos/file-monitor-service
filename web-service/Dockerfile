FROM golang:latest

RUN apt update \
    && apt install -y apt-utils sqlite3 git curl \
    && go get "github.com/mattn/go-sqlite3"

WORKDIR /home/webserver

COPY static_files/ .

COPY webserver.go .
RUN go build webserver.go

CMD ["./webserver"]

EXPOSE 8080/tcp