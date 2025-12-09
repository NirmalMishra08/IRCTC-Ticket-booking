CREATE TYPE user_role AS ENUM{
    'ADMIN',
    'USER'
}

CREATE TYPE provider as ENUM {
    'EMAIL',
    'APPLE',
    'PASSWORD'

}
CREATE TABLE user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
    fullname TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    roles user_role NOT NULL DEFAULT 'USER',
    password_hash TEXT ,
    provider provider NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);