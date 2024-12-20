FROM postgres:16.1
# Environment variables required for database creation
ENV POSTGRES_DB=shopify_db_dev
ENV POSTGRES_USER=shopify_user_dev
ENV POSTGRES_PASSWORD=shopify_password_dev

COPY ./shopify/up.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]
