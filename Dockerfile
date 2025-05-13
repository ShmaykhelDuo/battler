FROM debian:bookworm AS build

RUN apt update && apt install -y wget

ENV TENSORFLOW_FILENAME=libtensorflow-cpu-linux-x86_64-2.15.0.tar.gz

RUN wget -q --no-check-certificate https://storage.googleapis.com/tensorflow/libtensorflow/${TENSORFLOW_FILENAME} \
    && tar -C /usr/local -xzf ${TENSORFLOW_FILENAME} \
    && ldconfig /usr/local/lib

RUN wget https://github.com/pressly/goose/releases/download/v3.24.3/goose_linux_x86_64 -O goose && chmod +x goose

COPY ml/models ml/models
COPY migrations migrations
COPY battler battler

ENTRYPOINT [ "sh", "-c", "./goose postgres $DB_CONN up -dir migrations && ./battler" ]
