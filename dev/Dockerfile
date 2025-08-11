FROM python:3.11-slim

RUN apt-get update && \
    apt-get install -y redis-server && \
    rm -rf /var/lib/apt/lists/*

RUN pip install --no-cache-dir mysql-connector-python redis

COPY ./scripts ./scripts

CMD ["redis-server", "--protected-mode", "no"]
