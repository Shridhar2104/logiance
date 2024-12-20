FROM postgres:14-alpine

# Copy initialization scripts
COPY up.sql /docker-entrypoint-initdb.d/01_init.sql

# Set permissions for initialization scripts
RUN chmod 755 /docker-entrypoint-initdb.d/01_init.sql