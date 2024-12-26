DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_user WHERE usename = 'shipment_user_dev') THEN
        CREATE USER shipment_user_dev WITH PASSWORD 'shipment_password_dev';
    END IF;
END
$$;

ALTER USER shipment_user_dev WITH SUPERUSER;

DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'shipment_db_dev') THEN
        CREATE DATABASE shipment_db_dev;
    END IF;
END
$$;

GRANT ALL PRIVILEGES ON DATABASE shipment_db_dev TO shipment_user_dev;
\c shipment_db_dev
GRANT ALL ON SCHEMA public TO shipment_user_dev;