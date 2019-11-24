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
RUN cd /home/file-monitor-service && go build webserver.go
RUN chmod u+x script.sh
CMD ["/bin/bash", "./script.sh"]
CMD ["./webserver"]
EXPOSE 8080/tcp
