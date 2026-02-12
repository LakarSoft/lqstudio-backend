-- +goose Up
-- +goose StatementBegin
INSERT INTO themes (code, name) VALUES
    ('A', 'Theme A'),
    ('B', 'Theme B');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM themes WHERE code IN ('A', 'B');
-- +goose StatementEnd
