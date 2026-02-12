-- +goose Up
-- Core booking system tables for 20-minute slot booking system

-- Drop old tables if they exist
DROP TABLE IF EXISTS booking_theme_slots CASCADE;
DROP TABLE IF EXISTS bookings CASCADE;
DROP TABLE IF EXISTS packages CASCADE;
DROP TABLE IF EXISTS themes CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Packages table (slot-based)
CREATE TABLE packages (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    offers TEXT,
    slot_count INTEGER NOT NULL CHECK (slot_count IN (1, 2, 3)),
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    discount DECIMAL(10,2) DEFAULT 0 CHECK (discount >= 0),
    final_price DECIMAL(10,2) GENERATED ALWAYS AS (price - discount) STORED,
    image_url VARCHAR(500),
    booking_level VARCHAR(20) NOT NULL CHECK (booking_level IN ('theme', 'studio')),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT valid_booking_level CHECK (
        (slot_count IN (1, 2) AND booking_level = 'theme') OR
        (slot_count = 3 AND booking_level = 'studio')
    )
);

-- Themes table
CREATE TABLE themes (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add-ons table
CREATE TABLE add_ons (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Users table (guest bookings)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bookings table
CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    package_id INTEGER NOT NULL REFERENCES packages(id),
    booking_date DATE NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status VARCHAR(30) DEFAULT 'pending_verification'
        CHECK (status IN ('draft', 'pending_verification', 'confirmed', 'rejected')),
    booking_level VARCHAR(20) NOT NULL CHECK (booking_level IN ('theme', 'studio')),
    package_amount DECIMAL(10,2) NOT NULL CHECK (package_amount >= 0),
    add_ons_amount DECIMAL(10,2) DEFAULT 0 CHECK (add_ons_amount >= 0),
    total_amount DECIMAL(10,2) NOT NULL CHECK (total_amount >= 0),
    payment_confirmed_at TIMESTAMP,
    admin_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Ensure start times align to 20-minute boundaries
    CONSTRAINT valid_time_slot CHECK (
        EXTRACT(MINUTE FROM start_time)::INTEGER % 20 = 0
        AND EXTRACT(SECOND FROM start_time) = 0
    ),

    -- Ensure end time is after start time
    CONSTRAINT valid_time_range CHECK (end_time > start_time)
);

-- Booking theme slots (source of truth for availability)
CREATE TABLE booking_theme_slots (
    id SERIAL PRIMARY KEY,
    booking_id INTEGER NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    theme_id INTEGER REFERENCES themes(id),
    slot_number INTEGER NOT NULL CHECK (slot_number BETWEEN 1 AND 3),
    slot_start_time TIMESTAMP NOT NULL,
    slot_end_time TIMESTAMP NOT NULL,

    -- Each slot must be exactly 20 minutes
    CONSTRAINT valid_slot_duration CHECK (
        EXTRACT(EPOCH FROM (slot_end_time - slot_start_time)) / 60 = 20
    ),

    -- Unique constraint to prevent duplicate slots
    CONSTRAINT unique_booking_slot UNIQUE (booking_id, slot_number)
);

-- Booking add-ons (line items)
CREATE TABLE booking_add_ons (
    id SERIAL PRIMARY KEY,
    booking_id INTEGER NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    add_on_id INTEGER NOT NULL REFERENCES add_ons(id),
    quantity INTEGER DEFAULT 1 CHECK (quantity > 0),
    unit_price DECIMAL(10,2) NOT NULL CHECK (unit_price >= 0),
    total_price DECIMAL(10,2) GENERATED ALWAYS AS (quantity * unit_price) STORED,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_bookings_date_status ON bookings(booking_date, status);
CREATE INDEX idx_bookings_user ON bookings(user_id);
CREATE INDEX idx_booking_theme_slots_time ON booking_theme_slots(slot_start_time, slot_end_time);
CREATE INDEX idx_booking_theme_slots_theme ON booking_theme_slots(theme_id);
CREATE INDEX idx_booking_theme_slots_booking ON booking_theme_slots(booking_id);

-- +goose Down
DROP TABLE IF EXISTS booking_add_ons;
DROP TABLE IF EXISTS booking_theme_slots;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS add_ons;
DROP TABLE IF EXISTS themes;
DROP TABLE IF EXISTS packages;
