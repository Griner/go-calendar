FROM go-calendar

EXPOSE 8889
WORKDIR /
CMD /app/go-calendar grpcServer --config /config/config.docker.yml --address 0.0.0.0 --port 8889
