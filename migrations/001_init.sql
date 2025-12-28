-- 001_init.sql
CREATE TABLE IF NOT EXISTS documents (
    id UUID PRIMARY KEY,
    hash TEXT NOT NULL UNIQUE,
    hash_algo TEXT NOT NULL,
    size BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sign_events (
    id UUID PRIMARY KEY,
    document_id UUID NOT NULL REFERENCES documents(id),
    signer_email TEXT NOT NULL,
    ip_address TEXT,
    user_agent TEXT,
    signed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
