FROM ubuntu:latest

RUN apt update \
    && apt install -y inotify-tools bc sqlite git curl 
    
WORKDIR /home/monitor

COPY monitor_service.sh .

CMD ["/bin/bash", "monitor_service.sh"]