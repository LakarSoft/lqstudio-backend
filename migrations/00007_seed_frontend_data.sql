-- +goose Up

-- Seed admin user (password: Not2tell!)
-- bcrypt hash of "Not2tell!"
INSERT INTO users (id, email, password_hash, name, role) VALUES
('usr-admin-001', 'syedizzuddin@lakarsoft.com', '$2a$10$BTwHyEMrtqY7Y0A/R0cf8.RDwhkRnjVBSrGH0PxRvP.JXMNmkk6A2', 'Admin User', 'ADMIN')
ON CONFLICT (id) DO UPDATE SET password_hash = EXCLUDED.password_hash;

-- Seed packages (matching frontend exactly)
INSERT INTO packages (id, name, description, price, discount, offers, image_url) VALUES
('pkg-early-bird', 'Early Bird Promotion', 'Special rate for early bookings', 100.00, 10.00,
 '["1 Photoshoot Slot (15 min)", "Unlimited Soft Copy (Google Drive)", "Valid for 2-10 pax", "RM 10 additional per pax (above 6yo)"]',
 'https://images.unsplash.com/photo-1542038784456-1ea8e935640e'),
('pkg-family-frame', 'Family Frame Package', 'Includes a high-quality framed photo', 180.00, 0.00,
 '["1 Photoshoot Slot (15 min)", "Unlimited Soft Copy (Google Drive)", "Photo Frame (12 x 18 inch)", "Valid for 2-7 pax"]',
 'https://images.unsplash.com/photo-1511895426328-dc8714191300'),
('pkg-premium-hour', 'Premium 1-Hour Session', 'Our most comprehensive photography experience', 300.00, 0.00,
 '["3 Photoshoot Slots (1 Hour)", "Unlimited Soft Copy (Google Drive)", "Accommodates up to 15 pax", "RM 10 additional per pax"]',
 'https://images.unsplash.com/photo-1606103836293-0a063ee20566'),
('pkg-executive-grand', 'Executive Grand Frame', 'The ultimate package with large frame', 350.00, 0.00,
 '["1 Photoshoot Slot (15 min)", "Unlimited Soft Copy (Google Drive)", "Large Photo Frame (20 x 30 inch)", "RM 10 additional per pax"]',
 'https://images.unsplash.com/photo-1582555172866-f73bb12a2ab3');

-- Seed themes (matching frontend exactly)
INSERT INTO themes (id, name, description, image_url, price) VALUES
('theme-minimalist', 'Minimalist White', 'Clean, bright, and timeless aesthetic with neutral tones',
 'https://images.unsplash.com/photo-1494438639946-1ebd1d20bf85', 0.00),
('theme-vintage', 'Vintage Warmth', 'Cozy tones and nostalgic vibes with retro props',
 'https://images.unsplash.com/photo-1519741497674-611481863552', 0.00),
('theme-midnight', 'Midnight Elegance', 'Moody and sophisticated dark-themed setup',
 'https://images.unsplash.com/photo-1554048612-b6a482bc67e5', 0.00);

-- Seed addons (matching frontend exactly)
INSERT INTO addons (id, name, description, price, unit) VALUES
('addon-extra-pax', 'Extra Pax', 'Additional person in the studio (above 6 years old).', 10.00, 'pax'),
('addon-fast-delivery', 'Fast Delivery', 'Receive all soft copies within 24 hours.', 20.00, NULL),
('addon-physical-print', '4R Glossy Print', 'High-quality physical print of your choice.', 5.00, 'pc'),
('addon-makeup', 'Basic Makeup', 'Simple touch-up by our in-house artist.', 50.00, 'person');

-- +goose Down
DELETE FROM booking_addons;
DELETE FROM booking_slots;
DELETE FROM bookings;
DELETE FROM addons;
DELETE FROM themes;
DELETE FROM packages;
DELETE FROM users WHERE role = 'ADMIN';
