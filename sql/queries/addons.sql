-- name: GetAddonByID :one
SELECT * FROM addons
WHERE id = $1::varchar LIMIT 1;

-- name: GetAddonsByIDs :many
SELECT * FROM addons
WHERE id = ANY($1::varchar[]);

-- name: GetActiveAddons :many
SELECT * FROM addons
WHERE is_active = true
ORDER BY price ASC;

-- name: ListAllAddons :many
SELECT * FROM addons
ORDER BY name ASC;

-- name: CreateAddon :one
INSERT INTO addons (
    id, name, description, price, unit, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdateAddon :one
UPDATE addons
SET
    name = $2,
    description = $3,
    price = $4,
    unit = $5,
    is_active = $6,
    updated_at = NOW()
WHERE id = $1::varchar
RETURNING *;

-- name: DeleteAddon :exec
DELETE FROM addons WHERE id = $1::varchar;

-- name: ToggleAddonActive :one
UPDATE addons
SET is_active = NOT is_active,
    updated_at = NOW()
WHERE id = $1::varchar
RETURNING *;
