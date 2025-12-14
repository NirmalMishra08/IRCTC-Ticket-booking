-- name: CreateBooking: one
INSERT INTO booking (userId, trainId, travelDate, status, holdToken)
VALUES ($1, $2, $3, 'PENDING', $4)
RETURNING *;

-- name: CreateBookingItem :one
INSERT INTO bookingItem (bookingId, seatId, trainScheduleId)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateBookingStatus :exec
UPDATE booking
SET status = $2
WHERE id = $1;

-- name: GetBookingByHoldToken :one
SELECT * FROM booking WHERE holdToken = $1;

-- name: GetBookedSeats :many
SELECT bi.seatId
FROM bookingItem bi
JOIN booking b on bi.bookingId = b.id
 WHERE b.trainId = $2 AND
  b.travelDate = $3
  AND b.status IN('PENDING','CONFIRMED');

-- name: GetBookingbyUserId :many
SELECT bi.* , b.*
FROM booking b 
JOIN bookingItem bi 
ON b.id = bi.bookingId
WHERE b.userId = $1;

-- name: GetSeatsByTrainId :many
