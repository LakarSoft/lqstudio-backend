-- +goose Up
ALTER TABLE packages
ADD COLUMN module VARCHAR(50);

UPDATE packages
SET module = 'raya'
WHERE module IS NULL;

ALTER TABLE packages
ALTER COLUMN module SET DEFAULT 'raya';

ALTER TABLE packages
ALTER COLUMN module SET NOT NULL;

ALTER TABLE themes
ADD COLUMN module VARCHAR(50);

UPDATE themes
SET module = 'raya'
WHERE module IS NULL;

ALTER TABLE themes
ALTER COLUMN module SET DEFAULT 'raya';

ALTER TABLE themes
ALTER COLUMN module SET NOT NULL;

ALTER TABLE addons
ADD COLUMN module VARCHAR(50);

UPDATE addons
SET module = 'raya'
WHERE module IS NULL;

ALTER TABLE addons
ALTER COLUMN module SET DEFAULT 'raya';

ALTER TABLE addons
ALTER COLUMN module SET NOT NULL;

-- Seed convocation packages
INSERT INTO packages (id, module, name, description, duration_minutes, price, discount, offers, image_url, is_active) VALUES
('pkg-convo-basic', 'convocation', 'Convocation Basic', 'Single graduate studio session for a focused convocation portrait.', 30, 120.00, 0.00,
 '["1 Photoshoot Slot (30 min)", "1 Convocation Theme", "Unlimited Soft Copy (Google Drive)", "Valid for 1 graduate"]',
 'https://images.unsplash.com/photo-1523580846011-d3a5bc25702b?q=80&w=2340&auto=format&fit=crop', true),
('pkg-convo-family', 'convocation', 'Convocation Family', 'Convocation session for one graduate with family portraits.', 60, 220.00, 0.00,
 '["2 Photoshoot Slots (60 min)", "1 Convocation Theme", "Unlimited Soft Copy (Google Drive)", "Valid for 1 graduate and family"]',
 'https://images.unsplash.com/photo-1523580846011-d3a5bc25702b?q=80&w=2340&auto=format&fit=crop', true),
('pkg-convo-premium', 'convocation', 'Convocation Premium', 'Extended convocation session with more portrait variety and keepsake coverage.', 90, 320.00, 0.00,
 '["3 Photoshoot Slots (90 min)", "1 Convocation Theme", "Unlimited Soft Copy (Google Drive)", "Premium portrait selection"]',
 'https://images.unsplash.com/photo-1523580846011-d3a5bc25702b?q=80&w=2340&auto=format&fit=crop', true)
ON CONFLICT (id) DO UPDATE SET
    module = EXCLUDED.module,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    duration_minutes = EXCLUDED.duration_minutes,
    price = EXCLUDED.price,
    discount = EXCLUDED.discount,
    offers = EXCLUDED.offers,
    image_url = EXCLUDED.image_url,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- Seed convocation themes
INSERT INTO themes (id, module, name, description, image_url, price, is_active) VALUES
('theme-convo-formal', 'convocation', 'Formal Convocation', 'Classic academic portrait setup for robe and scroll shots.',
 'https://images.unsplash.com/photo-1523050854058-8df90110c9f1?q=80&w=2340&auto=format&fit=crop', 0.00, true),
('theme-convo-family', 'convocation', 'Family Convocation', 'Warm family-focused convocation backdrop for group portraits.',
 'https://images.unsplash.com/photo-1523050854058-8df90110c9f1?q=80&w=2340&auto=format&fit=crop', 0.00, true),
('theme-convo-premium', 'convocation', 'Premium Convocation', 'Elevated convocation setup for polished graduate portraits.',
 'https://images.unsplash.com/photo-1523050854058-8df90110c9f1?q=80&w=2340&auto=format&fit=crop', 0.00, true)
ON CONFLICT (id) DO UPDATE SET
    module = EXCLUDED.module,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    image_url = EXCLUDED.image_url,
    price = EXCLUDED.price,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- Seed convocation addons
INSERT INTO addons (id, module, name, description, price, unit, is_active) VALUES
('addon-convo-extra-graduate', 'convocation', 'Extra Graduate', 'Additional graduate included in the convocation session.', 50.00, 'pax', true),
('addon-convo-makeup', 'convocation', 'Makeup Touch-Up', 'Basic makeup and grooming touch-up before the session.', 80.00, 'person', true),
('addon-convo-premium-frame', 'convocation', 'Premium Convocation Frame', 'Printed convocation portrait with premium frame.', 150.00, 'pc', true)
ON CONFLICT (id) DO UPDATE SET
    module = EXCLUDED.module,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    price = EXCLUDED.price,
    unit = EXCLUDED.unit,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- +goose Down
ALTER TABLE packages
DROP COLUMN IF EXISTS module;

ALTER TABLE themes
DROP COLUMN IF EXISTS module;

ALTER TABLE addons
DROP COLUMN IF EXISTS module;
