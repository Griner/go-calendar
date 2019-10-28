FROM go-calendar

WORKDIR /
CMD /app/go-calendar eventNotifier --config /config/config.docker.yml
