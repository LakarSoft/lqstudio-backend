-- This file contains the database schema for sqlc to understand the structure
-- It should match the migrations but is used by sqlc for code generation
-- This matches migration 00006_convert_to_string_ids.sql

-- Users table (serves as both customer and admin)
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    role VARCHAR(50) NOT NULL DEFAULT 'CUSTOMER',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_users_email ON users(email);

-- Packages table
CREATE TABLE packages (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    duration_minutes INTEGER NOT NULL DEFAULT 20,
    price DECIMAL(10, 2) NOT NULL,
    discount DECIMAL(5, 2) DEFAULT 0,
    offers JSONB,
    image_url TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Themes table
CREATE TABLE themes (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    image_url TEXT NOT NULL DEFAULT '',
    price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Addons table
CREATE TABLE addons (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Bookings table
CREATE TABLE bookings (
    id VARCHAR(255) PRIMARY KEY,
    package_id VARCHAR(255) NOT NULL REFERENCES packages(id) ON DELETE RESTRICT,
    customer_name VARCHAR(255) NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(50) NOT NULL,
    customer_notes TEXT,
    payment_screenshot_url TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    total_price DECIMAL(10, 2) NOT NULL,
    admin_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_booking_status CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED'))
);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_created_at ON bookings(created_at DESC);
CREATE INDEX idx_bookings_email ON bookings(customer_email);

-- Booking Slots table
CREATE TABLE booking_slots (
    id SERIAL PRIMARY KEY,
    booking_id VARCHAR(255) NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    time TIME NOT NULL,
    theme_id VARCHAR(255) NOT NULL REFERENCES themes(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_booking_slots_booking_id ON booking_slots(booking_id);
CREATE INDEX idx_booking_slots_availability ON booking_slots(date, time, theme_id);

-- Booking Addons table
CREATE TABLE booking_addons (
    id SERIAL PRIMARY KEY,
    booking_id VARCHAR(255) NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    addon_id VARCHAR(255) NOT NULL REFERENCES addons(id) ON DELETE RESTRICT,
    quantity INT NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_booking_addons_booking_id ON booking_addons(booking_id);
