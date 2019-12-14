FROM golang:latest

RUN apt update \
    && apt install -y apt-utils sqlite3 git curl \
    && go get "github.com/mattn/go-sqlite3" \
    && groupadd -g 1000 appuser \
    && useradd -r -u 1000 -g appuser appuser

USER appuser

WORKDIR /home/webserver

RUN mkdir -p monitored_folder

COPY static_files/ .
COPY webserver.go .
RUN go build webserver.go

CMD ["./webserver"]

EXPOSE 8080/tcp
