FROM alpine:latest

WORKDIR /app
ADD ./convert /app/convert
ADD ./etc /etc/convert

EXPOSE 8085

ENTRYPOINT ["/app/convert", "-config", "/etc/convert/config.yaml"]
