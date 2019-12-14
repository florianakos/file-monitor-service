FROM golang:latest

RUN apt update \
    && apt install -y apt-utils inotify-tools bc sqlite3 git curl \
    && go get "github.com/mattn/go-sqlite3" \
    && mkdir -p /home/file-monitor-service \
    && groupadd -g 1000 appuser \
    && useradd -r -u 1000 -g appuser appuser

WORKDIR /home/webserver

COPY static_files/ .
COPY webserver.go .
RUN go build webserver.go

USER appuser
CMD ["./webserver"]

EXPOSE 8080/tcp