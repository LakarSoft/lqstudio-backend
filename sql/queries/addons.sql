-- name: GetAddonByID :one
SELECT * FROM addons
WHERE id = $1::varchar LIMIT 1;

-- name: GetAddonsByIDs :many
SELECT * FROM addons
WHERE id = ANY($1::varchar[]);

-- name: GetActiveAddons :many
SELECT * FROM addons
WHERE is_active = true
  AND ($1::text = '' OR module = $1)
ORDER BY module ASC, price ASC;

-- name: ListAllAddons :many
SELECT * FROM addons
WHERE ($1::text = '' OR module = $1)
ORDER BY module ASC, created_at DESC;

-- name: CreateAddon :one
INSERT INTO addons (
    id, module, name, description, price, unit, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: UpdateAddon :one
UPDATE addons
SET
    module = $2,
    name = $3,
    description = $4,
    price = $5,
    unit = $6,
    is_active = $7,
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
