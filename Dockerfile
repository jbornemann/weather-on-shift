FROM alpine:latest

COPY weather-on-shift /bin

ENTRYPOINT /bin/weather-on-shift
