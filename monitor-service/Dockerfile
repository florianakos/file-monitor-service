FROM ubuntu:latest

RUN apt update \
    && apt install -y inotify-tools bc sqlite git curl \
    && mkdir -p /home/file-monitor-service \
    && groupadd -g 1000 appuser \
    && useradd -r -u 1000 -g appuser appuser

WORKDIR /home/monitor

COPY monitor_service.sh .

USER appuser 

CMD ["/bin/bash", "monitor_service.sh"]