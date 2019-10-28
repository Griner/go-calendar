FROM go-calendar

WORKDIR /
CMD /app/go-calendar eventNotifierWorker --config /config/config.docker.yml
