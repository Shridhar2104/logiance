# ./shipment/db.dockerfile
FROM postgres:14-alpine

# Add any custom configuration if needed
COPY ./shipment/init.sql /docker-entrypoint-initdb.d/

# Set default configuration
ENV POSTGRES_HOST_AUTH_METHOD=trust