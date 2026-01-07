
-- name: GetPaymentAndTrain :one
SELECT p.amount,p.status as payment_status,b.travelDate , b.holdToken , b.status as booking_status FROM
booking b JOIN
payment p ON b.id = p.bookingId
WHERE b.trainId = $2 AND b.userId = $1 AND b.travelDate = $3
LIMIT 1;
