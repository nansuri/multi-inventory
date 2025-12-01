-- Multi Inventory Management System - Production Schema
-- Execute this script manually in production database
-- Run as a user with sufficient privileges (e.g., postgres superuser)

-- =============================================================================
-- STEP 1: Enable Required Extensions
-- =============================================================================

-- Enable UUID generation extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Enable UUID-OSSP extension (alternative UUID generation)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =============================================================================
-- STEP 2: Create Tables
-- =============================================================================

-- Users Table
-- Stores user authentication and profile information
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create unique index on username
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Create index on role for filtering
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Items Table
-- Stores inventory item information
CREATE TABLE IF NOT EXISTS items (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    barcode TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    location TEXT,
    is_halal BOOLEAN NOT NULL DEFAULT TRUE,
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create unique index on barcode
CREATE UNIQUE INDEX IF NOT EXISTS idx_items_barcode ON items(barcode);

-- Create index on name for searching
CREATE INDEX IF NOT EXISTS idx_items_name ON items(name);

-- Create index on is_halal for filtering
CREATE INDEX IF NOT EXISTS idx_items_is_halal ON items(is_halal);

-- Sales Orders Table
-- Stores sales order header information
CREATE TABLE IF NOT EXISTS sales_orders (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID,
    total_price DECIMAL(10, 2) NOT NULL CHECK (total_price >= 0),
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_sales_orders_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Create index on user_id for lookups
CREATE INDEX IF NOT EXISTS idx_sales_orders_user_id ON sales_orders(user_id);

-- Create index on status for filtering
CREATE INDEX IF NOT EXISTS idx_sales_orders_status ON sales_orders(status);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_sales_orders_created_at ON sales_orders(created_at DESC);

-- Sales Order Items Table
-- Stores individual items in each sales order
CREATE TABLE IF NOT EXISTS sales_order_items (
    id BIGSERIAL PRIMARY KEY,
    sales_order_id BIGINT NOT NULL,
    item_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price_at_sale DECIMAL(10, 2) NOT NULL CHECK (price_at_sale >= 0),
    is_fulfilled BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_sales_order_items_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_sales_order_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE RESTRICT
);

-- Create index on sales_order_id for lookups
CREATE INDEX IF NOT EXISTS idx_sales_order_items_order_id ON sales_order_items(sales_order_id);

-- Create index on item_id for lookups
CREATE INDEX IF NOT EXISTS idx_sales_order_items_item_id ON sales_order_items(item_id);

-- Create index on is_fulfilled for filtering
CREATE INDEX IF NOT EXISTS idx_sales_order_items_fulfilled ON sales_order_items(is_fulfilled);

-- =============================================================================
-- STEP 3: Add Comments (Documentation)
-- =============================================================================

COMMENT ON TABLE users IS 'User accounts with authentication information';
COMMENT ON COLUMN users.id IS 'Unique user identifier (UUID)';
COMMENT ON COLUMN users.username IS 'Unique username for login';
COMMENT ON COLUMN users.password IS 'Bcrypt hashed password';
COMMENT ON COLUMN users.role IS 'User role: admin, manager, user, etc.';

COMMENT ON TABLE items IS 'Inventory items with details';
COMMENT ON COLUMN items.id IS 'Unique item identifier';
COMMENT ON COLUMN items.barcode IS 'Unique barcode/QR code';
COMMENT ON COLUMN items.price IS 'Current selling price';
COMMENT ON COLUMN items.location IS 'Storage location in warehouse';
COMMENT ON COLUMN items.is_halal IS 'Halal certification status';
COMMENT ON COLUMN items.quantity IS 'Current stock quantity';

COMMENT ON TABLE sales_orders IS 'Sales order headers';
COMMENT ON COLUMN sales_orders.user_id IS 'User who created the order';
COMMENT ON COLUMN sales_orders.status IS 'Order status: pending, completed, cancelled';

COMMENT ON TABLE sales_order_items IS 'Line items within sales orders';
COMMENT ON COLUMN sales_order_items.price_at_sale IS 'Price at time of sale (historical)';
COMMENT ON COLUMN sales_order_items.is_fulfilled IS 'Whether item has been picked/fulfilled';

-- =============================================================================
-- STEP 4: Create Triggers for Updated At
-- =============================================================================

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for users table
DROP TRIGGER IF EXISTS trigger_users_updated_at ON users;
CREATE TRIGGER trigger_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for items table
DROP TRIGGER IF EXISTS trigger_items_updated_at ON items;
CREATE TRIGGER trigger_items_updated_at
    BEFORE UPDATE ON items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for sales_orders table
DROP TRIGGER IF EXISTS trigger_sales_orders_updated_at ON sales_orders;
CREATE TRIGGER trigger_sales_orders_updated_at
    BEFORE UPDATE ON sales_orders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =============================================================================
-- STEP 5: Insert Sample Data (Optional - for testing)
-- =============================================================================

-- Uncomment below to insert sample data for testing

-- Insert admin user (password is 'admin123' hashed with bcrypt)
-- INSERT INTO users (username, password, role) VALUES 
-- ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin')
-- ON CONFLICT (username) DO NOTHING;

-- Insert sample items
-- INSERT INTO items (name, barcode, price, location, is_halal, quantity) VALUES
-- ('Sample Item 1', '1234567890123', 10.50, 'A1', TRUE, 100),
-- ('Sample Item 2', '9876543210987', 25.00, 'B2', TRUE, 50)
-- ON CONFLICT (barcode) DO NOTHING;

-- =============================================================================
-- STEP 6: Grant Permissions (Adjust as needed)
-- =============================================================================

-- Grant permissions to application user (replace 'app_user' with your actual username)
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO app_user;

-- =============================================================================
-- VERIFICATION QUERIES
-- =============================================================================

-- Check if all tables are created
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
  AND table_name IN ('users', 'items', 'sales_orders', 'sales_order_items')
ORDER BY table_name;

-- Check if all indexes are created
SELECT tablename, indexname 
FROM pg_indexes 
WHERE schemaname = 'public' 
  AND tablename IN ('users', 'items', 'sales_orders', 'sales_order_items')
ORDER BY tablename, indexname;

-- Check if extensions are enabled
SELECT extname, extversion 
FROM pg_extension 
WHERE extname IN ('pgcrypto', 'uuid-ossp');

-- Check if triggers are created
SELECT trigger_name, event_manipulation, event_object_table
FROM information_schema.triggers
WHERE trigger_schema = 'public'
ORDER BY event_object_table, trigger_name;

-- =============================================================================
-- ROLLBACK SCRIPT (Use with caution - this will delete all data!)
-- =============================================================================

-- Uncomment below to drop all tables and start fresh
-- DROP TABLE IF EXISTS sales_order_items CASCADE;
-- DROP TABLE IF EXISTS sales_orders CASCADE;
-- DROP TABLE IF EXISTS items CASCADE;
-- DROP TABLE IF EXISTS users CASCADE;
-- DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;
