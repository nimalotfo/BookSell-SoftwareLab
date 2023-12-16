CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    mobile_number VARCHAR(12) NOT NULL UNIQUE,
    email VARCHAR(50),
    city VARCHAR(50)
);

CREATE INDEX idx_user_username ON users(username);