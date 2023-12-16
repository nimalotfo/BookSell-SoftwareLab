CREATE USER users_admin WITH ENCRYPTED PASSWORD 'pass1234';
CREATE DATABASE users_db;
GRANT ALL PRIVILEGES ON DATABASE users_db TO users_admin;
\connect users_db;
CREATE SCHEMA IF NOT EXISTS users_schema AUTHORIZATION users_admin;
ALTER DATABASE users_db SET search_path TO users_schema;
