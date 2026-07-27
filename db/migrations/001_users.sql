-- users 表
CREATE TABLE IF NOT EXISTS users
(
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    password_hash VARCHAR(100) NOT NULL,
    email         VARCHAR(255),

    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,

    created_at    timestamptz  NOT NULL DEFAULT now(),
    updated_at    timestamptz  NOT NULL DEFAULT now()
);

INSERT INTO users (username, password_hash, email)
VALUES ('test_user', '123456', 'example@eg.com');
