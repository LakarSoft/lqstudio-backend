-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

CREATE TABLE IF NOT EXISTS themes (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(1) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS packages (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(1) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    total_duration_minutes INT NOT NULL,
    is_exclusive BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS bookings (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    package_id BIGINT NOT NULL REFERENCES packages(id) ON DELETE RESTRICT,
    booking_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'CONFIRMED', 'CANCELLED')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_date ON bookings(booking_date);
CREATE INDEX idx_bookings_status ON bookings(status);

CREATE TABLE IF NOT EXISTS booking_theme_slots (
    id BIGSERIAL PRIMARY KEY,
    booking_id BIGINT NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    theme_id BIGINT NOT NULL REFERENCES themes(id) ON DELETE RESTRICT,
    slot_start_time TIMESTAMP NOT NULL,
    slot_end_time TIMESTAMP NOT NULL
);

CREATE INDEX idx_booking_theme_slots_booking_id ON booking_theme_slots(booking_id);
CREATE INDEX idx_booking_theme_slots_theme_date ON booking_theme_slots(theme_id, slot_start_time, slot_end_time);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS booking_theme_slots CASCADE;
DROP TABLE IF EXISTS bookings CASCADE;
DROP TABLE IF EXISTS packages CASCADE;
DROP TABLE IF EXISTS themes CASCADE;
DROP TABLE IF EXISTS users CASCADE;
-- +goose StatementEnd
