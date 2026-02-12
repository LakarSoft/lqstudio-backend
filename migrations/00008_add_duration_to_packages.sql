-- +goose Up
-- Add duration_minutes column to packages table
ALTER TABLE packages ADD COLUMN IF NOT EXISTS duration_minutes INTEGER NOT NULL DEFAULT 20;

-- Update existing packages with appropriate durations based on their descriptions
-- Early Bird: 1 slot (15-20 min) = 20 minutes
UPDATE packages SET duration_minutes = 20 WHERE id = 'pkg-early-bird';

-- Family Frame: 1 slot (15-20 min) = 20 minutes
UPDATE packages SET duration_minutes = 20 WHERE id = 'pkg-family-frame';

-- Premium 1-Hour: 3 slots = 60 minutes
UPDATE packages SET duration_minutes = 60 WHERE id = 'pkg-premium-hour';

-- Executive Grand: 1 slot (15-20 min) = 20 minutes
UPDATE packages SET duration_minutes = 20 WHERE id = 'pkg-executive-grand';

-- +goose Down
ALTER TABLE packages DROP COLUMN IF EXISTS duration_minutes;
