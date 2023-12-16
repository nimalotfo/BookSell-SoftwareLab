CREATE USER review_admin WITH ENCRYPTED PASSWORD 'pass1234';
CREATE DATABASE review_db;
GRANT ALL PRIVILEGES ON DATABASE review_db TO review_admin;
\connect review_db;
CREATE SCHEMA IF NOT EXISTS review_schema AUTHORIZATION review_admin;
ALTER DATABASE review_db SET search_path TO review_schema;

