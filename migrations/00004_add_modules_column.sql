-- +goose Up
ALTER TABLE packages
ADD COLUMN module VARCHAR(50);

UPDATE packages
SET module = 'raya'
WHERE module IS NULL;

ALTER TABLE packages
ALTER COLUMN module SET NOT NULL;

ALTER TABLE themes
ADD COLUMN module VARCHAR(50);

UPDATE themes
SET module = 'raya'
WHERE module IS NULL;

ALTER TABLE themes
ALTER COLUMN module SET NOT NULL;

ALTER TABLE addons
ADD COLUMN module VARCHAR(50);

UPDATE addons
SET module = 'raya'
WHERE module IS NULL;

ALTER TABLE addons
ALTER COLUMN module SET NOT NULL;

-- +goose Down
ALTER TABLE packages
DROP COLUMN IF EXISTS module;

ALTER TABLE themes
DROP COLUMN IF EXISTS module;

ALTER TABLE addons
DROP COLUMN IF EXISTS module;
