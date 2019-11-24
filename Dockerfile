FROM ubuntu:latest


RUN apt update && apt install -y \
    inotify-tools \
    golang \
    bc  \
    sqlite3 \
    git \
    curl \

RUN go get github.com/mattn/go-sqlite3


RUN mkdir -p /home/file-monitor-service

WORKDIR /home/file-monitor-service

COPY . .

RUN go build webserver.go

CMD ["./script.sh"]
CMD ["./webserver"]

EXPOSE 8080/tcp
