CREATE USER book_admin WITH ENCRYPTED PASSWORD 'pass1234';
CREATE DATABASE book_db;
GRANT ALL PRIVILEGES ON DATABASE book_db TO book_admin;
\connect book_db;
CREATE SCHEMA IF NOT EXISTS book_schema AUTHORIZATION book_admin;
ALTER DATABASE book_db SET search_path TO book_schema;

