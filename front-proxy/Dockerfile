FROM golang:latest

COPY main.go /proxy_main.go
COPY entrypoint /entrypoint

RUN chmod +x /entrypoint

ENTRYPOINT ["/entrypoint"]
