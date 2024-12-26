FROM postgres:16.1
COPY ./shipment/init.sql /docker-entrypoint-initdb.d/