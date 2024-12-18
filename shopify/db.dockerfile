FROM postgres:16.1

COPY ./shopify/up.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]
