-- name: CreateBooking :one
INSERT INTO booking (userId, journey_id, booking_type, status, holdToken)
VALUES ($1, $2, 'NORMAL', 'PENDING', $3)
RETURNING *;

-- name: CreateBookingItem :one
INSERT INTO bookingItem (bookingId, seatId)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateBookingStatus :exec
UPDATE booking
SET status = $2
WHERE id = $1;

-- name: GetBookingByHoldToken :one
SELECT * FROM booking WHERE holdToken = $1;

-- name: GetBookedSeats :many
SELECT si.seat_id
FROM seat_inventory si
JOIN booking b ON si.booking_id = b.id
WHERE b.journey_id = $1
  AND si.status IN ('HELD','CONFIRMED');


-- name: GetBookingbyUserId :many
SELECT bi.* , b.*
FROM booking b 
JOIN bookingItem bi 
ON b.id = bi.bookingId
WHERE b.userId = $1;

-- name: GetActiveBookingByUser :one
SELECT *
FROM booking
WHERE userid = $1
  AND status = 'PENDING'
  AND createdAt > now() - INTERVAL '10 minutes'
ORDER BY createdAt DESC
LIMIT 1;

-- name: ExpireOldBooking :exec
UPDATE booking
SET status='EXPIRED'
 WHERE status='PENDING'
   AND createdat < now() - INTERVAL '10 minutes';


-- name: GetBookingItemsByBooking :many
SELECT seatId FROM bookingItem WHERE bookingId = $1;


-- name: DeleteBookingItemsByBooking :exec
DELETE FROM bookingItem WHERE bookingId = $1;


-- name: CountActiveBookingByTrain :one
SELECT COUNT(*)
FROM booking
WHERE journey_id = $1
  AND status = 'PENDING';


-- name: CurrentAvailableSeats :many
SELECT seat_id
FROM seat_inventory
WHERE journey_id = $1
  AND status = 'AVAILABLE'
  AND seat_id = ANY($2::int[]);



-- name: CreatePayment :one
 INSERT into payment (bookingId,amount,transactionId)
 VALUES($1,$2,$3)
 RETURNING *;

-- name: UpdateBookingItemStatus :exec
UPDATE bookingItem SET bookingStatus = $2 WHERE bookingId = $1;

-- name: UpdatePaymentStatus :exec
UPDATE payment SET status = $2 WHERE bookingId = $1;

-- name: GetBookingLockContext :many
SELECT
    b.journey_id,
    si.seat_id,
    b.holdToken
FROM booking b
JOIN seat_inventory si ON si.booking_id = b.id
WHERE b.id = $1;

-- name: DeleteBookingItem :exec
DELETE FROM
bookingItem b WHERE bookingId = $1 ;





