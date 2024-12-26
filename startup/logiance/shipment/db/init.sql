CREATE USER shipment_user_dev WITH PASSWORD 'shipment_password_dev';
ALTER USER shipment_user_dev WITH SUPERUSER;
CREATE DATABASE shipment_db_dev;
GRANT ALL PRIVILEGES ON DATABASE shipment_db_dev TO shipment_user_dev;
\c shipment_db_dev
GRANT ALL ON SCHEMA public TO shipment_user_dev;
