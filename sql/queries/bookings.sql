-- =====================================================
-- BOOKING QUERIES
-- =====================================================

-- name: CreateBooking :one
INSERT INTO bookings (
    id, package_id, customer_name, customer_email, customer_phone,
    customer_notes, status, total_price
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetBookingByID :one
SELECT * FROM bookings
WHERE id = $1::varchar LIMIT 1;

-- name: ListBookings :many
SELECT * FROM bookings
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountBookings :one
SELECT COUNT(*) FROM bookings;

-- name: GetBookingsByStatus :many
SELECT * FROM bookings
WHERE status = $1
ORDER BY created_at DESC;

-- name: GetBookingsByCustomerEmail :many
SELECT * FROM bookings
WHERE customer_email = $1
ORDER BY created_at DESC;

-- name: UpdateBookingStatus :one
UPDATE bookings
SET
    status = $2,
    admin_notes = $3,
    updated_at = NOW()
WHERE id = $1::varchar
RETURNING *;

-- name: UpdateBookingPaymentScreenshot :one
UPDATE bookings
SET
    payment_screenshot_url = $2,
    updated_at = NOW()
WHERE id = $1::varchar
RETURNING *;

-- =====================================================
-- BOOKING SLOTS QUERIES
-- =====================================================

-- name: CreateBookingSlot :one
INSERT INTO booking_slots (
    booking_id, date, time, theme_id
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetBookingSlots :many
SELECT * FROM booking_slots
WHERE booking_id = $1::varchar
ORDER BY date, time;

-- name: GetBookedSlotsForThemeAndDate :many
SELECT bs.* FROM booking_slots bs
INNER JOIN bookings b ON bs.booking_id = b.id
WHERE bs.theme_id = $1::varchar
  AND bs.date = $2
  AND b.status IN ('PENDING', 'APPROVED')
ORDER BY bs.time;

-- name: GetBookedSlotsForAllThemesAndDate :many
SELECT bs.* FROM booking_slots bs
INNER JOIN bookings b ON bs.booking_id = b.id
WHERE bs.date = $1
  AND b.status IN ('PENDING', 'APPROVED')
ORDER BY bs.time;

-- =====================================================
-- BOOKING ADDONS QUERIES
-- =====================================================

-- name: CreateBookingAddon :one
INSERT INTO booking_addons (
    booking_id, addon_id, quantity
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetBookingAddons :many
SELECT * FROM booking_addons
WHERE booking_id = $1::varchar;
