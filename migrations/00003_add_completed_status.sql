-- +goose Up
-- +goose StatementBegin
ALTER TABLE bookings DROP CONSTRAINT valid_booking_status;
ALTER TABLE bookings ADD CONSTRAINT valid_booking_status CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED', 'COMPLETED'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE bookings SET status = 'APPROVED' WHERE status = 'COMPLETED';
ALTER TABLE bookings DROP CONSTRAINT valid_booking_status;
ALTER TABLE bookings ADD CONSTRAINT valid_booking_status CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED'));
-- +goose StatementEnd
