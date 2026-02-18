
-- name: GetPaymentAndTrain :one
SELECT b.* , p.*
FROM
booking b JOIN
payment p ON b.id = p.bookingId
WHERE b.journey_id = $2 AND b.userId = $1 AND p.status = "SUCCESS";


-- name: CreateRefund :one
INSERT INTO refund (userId , bookingId , amount , status , createdAt, updatedAt) 
VALUES ( $1 , $2 , $3 , $4 , now() , now() )
RETURNING *;

-- name: ReleaseSeatsByBooking :exec
UPDATE seat_inventory
SET status = 'AVAILABLE',
    booking_id = NULL
WHERE booking_id = $1
  AND status = 'CONFIRMED';

