-- +goose Up

-- Seed admin users
-- usr-admin-001 password: Not2tell!
-- usr-admin-002 password: Lqstudio123@
INSERT INTO users (id, email, password_hash, name, role) VALUES
('usr-admin-001', 'syedizzuddin@lakarsoft.com', '$2a$10$BTwHyEMrtqY7Y0A/R0cf8.RDwhkRnjVBSrGH0PxRvP.JXMNmkk6A2', 'Admin User', 'ADMIN'),
('usr-admin-002', 'leqiucreative@gmail.com', '$2a$10$aGc9.07eRjuQaeEf/Wc6qurTtvCjJq1Sb4UNWhf0dP5oe/BRat/8y', 'LQ Studio Admin', 'ADMIN')
ON CONFLICT (id) DO UPDATE SET password_hash = EXCLUDED.password_hash;

-- Seed packages
INSERT INTO packages (id, name, description, duration_minutes, price, discount, offers, image_url) VALUES
('pkg-essential', 'Essential Session', 'Perfect for quick portrait sessions', 20, 100.00, 10.00,
 '["1 Photoshoot Slot (20 min)", "Unlimited Soft Copy (Google Drive)", "Valid for 2-10 pax", "RM 10 additional per pax (above 6yo)"]',
 'https://images.unsplash.com/photo-1769283039142-586f352d126a?q=80&w=2340&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'),
('pkg-keepsake', 'Keepsake Session', 'Extended session with multiple themes', 40, 180.00, 0.00,
 '["1 Photoshoot Slot (40 min)", "Unlimited Soft Copy (Google Drive)", "Photo Frame (12 x 18 inch)", "Valid for 2-7 pax"]',
 'https://images.unsplash.com/photo-1769283039142-586f352d126a?q=80&w=2340&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'),
('pkg-group', 'Group Session', 'Our most comprehensive photography experience', 60, 300.00, 0.00,
 '["3 Photoshoot Slots (60 min)", "Unlimited Soft Copy (Google Drive)", "Accommodates up to 15 pax", "RM 10 additional per pax"]',
 'https://images.unsplash.com/photo-1769283039142-586f352d126a?q=80&w=2340&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'),
('pkg-classic', 'Classic Package', 'The ultimate package with frame', 20, 250.00, 0.00,
 '["1 Photoshoot Slot (20 min)", "Unlimited Soft Copy (Google Drive)", "Photo Frame (16 x 24 inch)", "RM 10 additional per pax"]',
 'https://images.unsplash.com/photo-1769283039142-586f352d126a?q=80&w=2340&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'),
('pkg-plus', 'Plus Package', 'The ultimate package with frames', 20, 300.00, 0.00,
 '["1 Photoshoot Slot (20 min)", "Unlimited Soft Copy (Google Drive)", "Photo Frame (12 x 18 inch)", "Photo Frame (16 x 24 inch)", "RM 10 additional per pax"]',
 'https://images.unsplash.com/photo-1769283039142-586f352d126a?q=80&w=2340&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'),
('pkg-signature', 'Signature Package', 'The ultimate package with large frames', 20, 350.00, 0.00,
 '["1 Photoshoot Slot (20 min)", "Unlimited Soft Copy (Google Drive)", "Photo Frame (12 x 18 inch)", "Photo Frame (20 x 30 inch)", "RM 10 additional per pax"]',
 'https://images.unsplash.com/photo-1769283039142-586f352d126a?q=80&w=2340&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D');

-- Seed themes
INSERT INTO themes (id, name, description, image_url, price) VALUES
('theme-A', 'Theme A', 'This is a placeholder for Theme A',
 'https://images.unsplash.com/photo-1768562821733-be6788db0666?q=80&w=2340&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D', 0.00),
('theme-B', 'Theme B', 'This is a placeholder for Theme B',
 'https://images.unsplash.com/photo-1768562821733-be6788db0666?q=80&w=2340&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D', 0.00);

-- Seed addons
INSERT INTO addons (id, name, description, price, unit) VALUES
('addon-frame-a', 'Photo Frame A', '12x18 inches', 90.00, 'pc'),
('addon-frame-b', 'Photo Frame B', '14x21 inches', 120.00, 'pc'),
('addon-frame-c', 'Photo Frame C', '16x24 inches', 150.00, 'pc'),
('addon-frame-d', 'Photo Frame D', '20x30 inches', 300.00, 'pc'),
('addon-family-book', 'Family Book', '48 pages', 100.00, 'pc'),
('addon-premium-book', 'Premium Book', '20 pages', 250.00, 'pc');

-- +goose Down
DELETE FROM booking_addons;
DELETE FROM booking_slots;
DELETE FROM bookings;
DELETE FROM addons;
DELETE FROM themes;
DELETE FROM packages;
DELETE FROM users WHERE role = 'ADMIN';
