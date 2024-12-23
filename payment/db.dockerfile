# Use the official PostgreSQL image
FROM postgres:16.1

# Set environment variables for default database, user, and password
ENV POSTGRES_DB=payment_service_db
ENV POSTGRES_USER=payment_service_user
ENV POSTGRES_PASSWORD=securepassword

# Expose the default PostgreSQL port
EXPOSE 5432

# Copy the SQL migration file to the Docker entrypoint directory
COPY ./payment/up.sql /docker-entrypoint-initdb.d/1.sql
