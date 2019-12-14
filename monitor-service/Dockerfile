FROM ubuntu:latest

RUN apt update \
    && apt install -y inotify-tools bc sqlite git curl \
    && groupadd -g 1000 appuser \
    && useradd -r -u 1000 -g appuser appuser

USER appuser

WORKDIR /home/monitor

RUN mkdir -p monitored_folder \
    && mkdir -p processed_files \
    && mkdir -p database \
    && mkdir -p logs

COPY monitor_service.sh .

CMD ["/bin/bash", "monitor_service.sh"]
