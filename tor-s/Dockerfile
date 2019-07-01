FROM ubuntu:xenial

ENV TOR_HOME /var/lib/tor
ENV TOR_PROXY_PORT 9050
ENV SERVICE_URL 127.0.0.1:80
### TOR
#RUN echo deb http://deb.torproject.org/torproject.org xenial main >> /etc/apt/sources.list.d/tor.list && \
#    echo deb-src http://deb.torproject.org/torproject.org xenial main >> /etc/apt/sources.list.d/tor.list && \
#    gpg --keyserver keys.gnupg.net --recv A3C4F0F979CAA22CDBA8F512EE8CBC9E886DDD89 && \
#    gpg --export A3C4F0F979CAA22CDBA8F512EE8CBC9E886DDD89 | apt-key add - && \
#    apt-get update && \
#    apt-get install -y tor deb.torproject.org-keyring gettext-base && \
#    rm -rf /var/lib/apt/lists/*
RUN apt update && apt install tor -y

RUN apt install polipo -y

COPY torrc.template /torrc.template
COPY entrypoint /entrypoint
COPY polipo.config /polipo.config
COPY polipo-run.sh /polipo-run.sh
RUN chmod +x /polipo-run.sh
RUN chmod +x /entrypoint && chmod u+r /torrc.template && \
    chown debian-tor:debian-tor /entrypoint /torrc.template

USER root
ENTRYPOINT ["/entrypoint"]
