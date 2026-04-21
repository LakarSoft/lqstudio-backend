-- +goose Up
ALTER TABLE packages
ADD COLUMN module VARCHAR(50);

UPDATE packages
SET module = 'raya'
WHERE module IS NULL;

ALTER TABLE packages
ALTER COLUMN module SET NOT NULL;

-- +goose Down
ALTER TABLE packages
DROP COLUMN IF EXISTS module;
