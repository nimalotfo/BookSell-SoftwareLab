CREATE TABLE roles(
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL
);

CREATE TABLE user_roles(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    role_id BIGINT REFERENCES roles(id),
    UNIQUE(user_id, role_id)
);

CREATE TABLE permissions(
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL
);

CREATE TABLE role_permissions(
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT REFERENCES roles(id),
    permission_id BIGINT REFERENCES permissions(id),
    UNIQUE(role_id, permission_id)
);

insert into roles values
	(1, 'admin'),
	(2, 'user')
;