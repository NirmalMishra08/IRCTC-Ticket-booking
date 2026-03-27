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

-- name: GetNextWaitlist :one
SELECT *
FROM waitlist
WHERE journey_id = $1
  AND status = 'WAITING'
ORDER BY priority_level DESC, waitlist_number ASC
LIMIT 1
FOR UPDATE SKIP LOCKED;

-- name: GetWaitlistBatch :many
SELECT *
FROM waitlist
WHERE journey_id = $1
  AND status = 'WAITING'
ORDER BY priority_level DESC, waitlist_number ASC
LIMIT $2;

-- name: CancelWaitlist :exec
UPDATE waitlist
SET status = 'CANCELLED',
    updatedAt = now()
WHERE bookingId = $1;

-- name: DeleteWaitlist :exec
DELETE FROM waitlist
WHERE bookingId = $1;

-- name: UpdateWaitlistStatus :exec
UPDATE waitlist
SET status = $2,
    updatedAt = now()
WHERE id = $1;