CREATE TABLE stores (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index on name
CREATE INDEX idx_stores_name ON stores(name);

-- Index on created_at
CREATE INDEX idx_stores_created_at ON stores(created_at);