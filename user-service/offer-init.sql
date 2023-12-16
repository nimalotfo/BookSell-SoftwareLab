CREATE USER offers_admin WITH ENCRYPTED PASSWORD 'pass1234';
CREATE DATABASE offers_db;
GRANT ALL PRIVILEGES ON DATABASE offers_db TO offers_admin;
\connect offers_db;
CREATE SCHEMA IF NOT EXISTS offers_schema AUTHORIZATION offers_admin;
ALTER DATABASE offers_db SET search_path TO offers_schema;

