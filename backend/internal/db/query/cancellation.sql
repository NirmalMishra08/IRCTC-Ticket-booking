
-- name: GetPaymentAndTrain :one
SELECT p.amount,p.status as payment_status,b.id as booking_id , b.travelDate , b.holdToken , b.status as booking_status FROM
booking b JOIN
payment p ON b.id = p.bookingId
WHERE b.trainId = $2 AND b.userId = $1 AND b.travelDate = $3
LIMIT 1;

-- name: CreateRefund :one
INSERT INTO refund (userId , bookingId , amount , status , createdAt, updatedAt) 
VALUES ( $1 , $2 , $3 , $4 , now() , now() )
RETURNING *;

