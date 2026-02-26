-- name: GetNextWaitlistNumber :one
select COALESCE(MAX(waitlist_number), 0) + 1
from waitlist
WHERE journey_id = $1
FOR UPDATE;

-- name: InsertWaitlist :exec
INSERT INTO waitlist (
  journey_id,
  bookingId,
  waitlist_number,
  status
)
VALUES ($1, $2, $3, 'WAITING');