FROM debian

RUN apt-get update && \
    apt-get install -y libmagickwand-dev
COPY assets /usr/local/share/liveserver

COPY ca-certificates.crt /etc/ssl/certs/
COPY bin/liveserver /
COPY db /db

EXPOSE 8080
EXPOSE 8081

ENTRYPOINT ["/liveserver"]
