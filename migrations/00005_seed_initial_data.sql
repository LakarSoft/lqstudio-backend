-- +goose Up
-- Seed initial themes
INSERT INTO themes (code, name, description, is_active) VALUES
    ('VINTAGE', 'Vintage Studio', 'Classic vintage photography setup with retro props', true),
    ('MODERN', 'Modern Minimalist', 'Clean modern aesthetic with neutral tones', true),
    ('GARDEN', 'Garden Paradise', 'Natural outdoor garden setting', true);

-- Seed initial packages
INSERT INTO packages (name, description, offers, slot_count, price, discount, booking_level, is_active) VALUES
    (
        '1 Slot Package',
        'Perfect for quick portrait sessions',
        '20 minutes, 1 theme of your choice, 10 edited photos',
        1,
        100.00,
        0,
        'theme',
        true
    ),
    (
        '2 Slots Package',
        'Extended session with multiple themes',
        '40 minutes, 2 themes (can be same theme), 20 edited photos',
        2,
        180.00,
        0,
        'theme',
        true
    ),
    (
        '3 Slots Package - Studio Exclusive',
        'Complete studio access for premium shoots',
        '60 minutes, all themes included, 30 edited photos, exclusive studio access',
        3,
        300.00,
        50,
        'studio',
        true
    );

-- Seed initial add-ons
INSERT INTO add_ons (name, description, price, is_active) VALUES
    ('Extra Photo Prints', '10 high-quality printed photos (6x8 inches)', 50.00, true),
    ('Digital Album', 'Custom-designed digital photo album', 80.00, true),
    ('Video Highlights', '2-minute highlight video of your session', 120.00, true),
    ('Professional Hair & Makeup', 'Professional styling before your shoot', 150.00, true);

-- +goose Down
DELETE FROM add_ons;
DELETE FROM packages;
DELETE FROM themes;
