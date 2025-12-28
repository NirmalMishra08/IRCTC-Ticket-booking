-- name: CreateBooking :one
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
 WHERE b.trainId = $1 AND
  b.travelDate = $2
  AND b.status IN('PENDING','CONFIRMED');

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
  AND createdat > now() - INTERVAL '10 minutes'
LIMIT 1;

-- name: ExpireOldBooking :exec
UPDATE booking
SET status='EXPIRED'
 WHERE status='PENDING'
   AND createdat > now() - INTERVAL '10 minutes';


-- name: GetBookingItemsByBooking :many
SELECT seatId FROM bookingItem WHERE bookingId = $1;


-- name: DeleteBookingItemsByBooking :exec
DELETE FROM bookingItem WHERE bookingId = $1;


-- name: CountActiveBookingByTrain :one
SELECT COUNT(*)
FROM booking
WHERE trainId = $1
     AND travelDate = $2
     AND status = 'PENDING';

-- name: CurrentAvailabeSeats :many
SELECT s.id
FROM seat s WHERE
s.id = ANY($1 :: int[])
 AND NOT EXISTS (
   SELECT 1 FROM
   bookingItem bi
   JOIN booking b ON bi.bookingId = b.id
   WHERE bi.seatId = s.id
   AND b.status  IN('PENDING','CONFIRMED')
   AND b.trainId = $2
   AND b.travelDate = $3
 );

 CREATE Table payment (
    id  SERIAL PRIMARY KEY ,
    bookingId  INTEGER,
    amount   FLOAT NOT NULL,
    status payment_status DEFAULT 'PENDING',
    transactionId TEXT NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT now()
);

-- name: CreatePayment :one
 INSERT into payment (bookingId,amount,transactionId)
 VALUES($1,$2,$3)
 RETURNING *;


