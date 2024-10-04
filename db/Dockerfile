FROM alpine:3.19.1

WORKDIR /migrations

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.19.2/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

ADD db/*.sql /migrations/
ADD db/migrate.sh /migrations/
RUN chmod +x /migrations/migrate.sh

ENTRYPOINT ["/migrations/migrate.sh"]
