CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       username VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       age INT,
                       gender BOOLEAN,
                       is_online BOOLEAN DEFAULT false
);

CREATE TABLE messages (
                          id SERIAL PRIMARY KEY,
                          sender_id UUID REFERENCES users(id),
                          receiver_id UUID REFERENCES users(id),
                          content TEXT NOT NULL,
                          seen BOOLEAN DEFAULT false
);
CREATE TABLE app_logs (
                          id SERIAL PRIMARY KEY,
                          level VARCHAR(10) NOT NULL,
                          message TEXT NOT NULL,
                          error_detail TEXT,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);