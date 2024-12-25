FROM postgres:16.1

COPY ./account/up.sql /docker-entrypoint-initdb.d/init.sql
RUN chmod 0755 /docker-entrypoint-initdb.d/init.sql

EXPOSE 5432

CMD ["postgres"]