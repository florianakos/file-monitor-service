FROM ubuntu:latest
RUN apt update && apt install -y \
    apt-utils \
    inotify-tools \
    golang \
    bc  \
    sqlite3 \
    git \
    curl
RUN go get "github.com/mattn/go-sqlite3"
RUN mkdir -p /home/file-monitor-service
WORKDIR /home/file-monitor-service
COPY . .
RUN go build webserver.go
CMD ["/bin/bash", "startup.sh"]
EXPOSE 8080/tcp
