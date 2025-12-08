
CREATE TYPE user_status AS ENUM ('http', 'ping', 'pending');
CREATE TYPE monitor_status AS ENUM ('up', 'down', 'unknown', 'pending');

CREATE INDEX idx_monitors_user_id ON monitors(user_id);
CREATE INDEX idx_monitor_logs_monitor_id ON monitor_logs(monitor_id);
CREATE INDEX idx_monitor_logs_checked_at ON monitor_logs(checked_at);


CREATE TYPE provider AS ENUM (
    'email',
    'google.com',
    'apple',
    'password'
);

CREATE TYPE user_role AS ENUM (
    'ADMIN',
    'USER'
);

CREATE TYPE profile_status AS ENUM (
    'PENDING',
    'COMPLETED'
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone TEXT,
    email TEXT UNIQUE NOT NULL,
    fullname TEXT NOT NULL,
    provider provider NOT NULL,
    password_hash TEXT,
    role user_role NOT NULL DEFAULT 'USER',
    profile_status profile_status DEFAULT 'PENDING',
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE TABLE user_profile (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    is_premium BOOLEAN DEFAULT FALSE,
    stripe_id VARCHAR(255),
    name VARCHAR(100),
    bio TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE monitors (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    method TEXT DEFAULT 'GET',
    type TEXT DEFAULT 'http',
    interval INTEGER NOT NULL,
    status monitor_status DEFAULT 'unknown',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE monitor_logs (
    id SERIAL PRIMARY KEY,
    monitor_id INTEGER REFERENCES monitors(id),
    status_code INTEGER,
    response_time FLOAT,
    dns_ok BOOLEAN,
    ssl_ok BOOLEAN,
    content_ok BOOLEAN,
    screenshot_url TEXT,
    checked_at TIMESTAMP DEFAULT now()
);


CREATE TABLE ALERT(
    id SERIAL PRIMARY KEY,
    monitor_id INTEGER REFERENCES monitors(id),
    alert_type TEXT,
    message TEXT,
    sent_at TIMESTAMP DEFAULT now(),
    created_at TIMESTAMP DEFAULT now()
);


CREATE TABLE subscriptions(
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    stripe_plan TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE analytics(
    id SERIAL PRIMARY KEY,
    monitor_id INTEGER REFERENCES monitors(id),
    uptime_percentage FLOAT,
    avg_response_time FLOAT,
    last_24h_downtime INTEGER,
    updated_at TIMESTAMP DEFAULT now()
);

-- Create indexes after all tables are created
CREATE INDEX idx_monitors_user_id ON monitors(user_id);
CREATE INDEX idx_monitor_logs_monitor_id ON monitor_logs(monitor_id);
CREATE INDEX idx_monitor_logs_checked_at ON monitor_logs(checked_at);