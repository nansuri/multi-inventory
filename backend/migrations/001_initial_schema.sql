-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS items (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    barcode TEXT NOT NULL UNIQUE,
    price DECIMAL(10, 2) NOT NULL,
    location TEXT,
    is_halal BOOLEAN DEFAULT TRUE,
    quantity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sales_orders (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    total_price DECIMAL(10, 2) NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending', -- pending, completed, cancelled
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sales_order_items (
    id BIGSERIAL PRIMARY KEY,
    sales_order_id BIGINT REFERENCES sales_orders(id) ON DELETE CASCADE,
    item_id BIGINT REFERENCES items(id),
    quantity INTEGER NOT NULL,
    price_at_sale DECIMAL(10, 2) NOT NULL,
    is_fulfilled BOOLEAN DEFAULT FALSE
);
