-- +goose Up
-- +goose StatementBegin
INSERT INTO packages (code, name, total_duration_minutes, is_exclusive) VALUES
    ('A', 'Package A - Single Theme', 30, FALSE),
    ('B', 'Package B - Sequential Themes', 60, FALSE),
    ('C', 'Package C - Exclusive Studio', 120, TRUE);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM packages WHERE code IN ('A', 'B', 'C');
-- +goose StatementEnd
