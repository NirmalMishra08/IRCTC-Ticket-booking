-- name: CreateTrain :one
INSERT into train(trainNumber,trainName,source,destination )
VALUES ( $1 ,$2 ,$3,$4) 
RETURNING *;

-- name: CreateTrainSchedule :one
INSERT into train_schedule (trainId,day,arrivalTime,departureTime)
VALUES ( $1 ,$2 ,$3,$4) 
RETURNING *;

-- name: GetTrainById :one
SELECT * FROM train
WHERE id = $1;
-- name: CreateCoach :one
INSERT into coach (trainId,coachtype,coachNumber) VALUES ($1 , $2 , $3) RETURNING *;

-- name: CreateSeat :one
INSERT into seat (coachId,seatno,berth) VALUES ($1 , $2 , $3) RETURNING *;

-- name: GetTrainScheduleByDay :one
SELECT *
FROM train_schedule
WHERE trainId = $1
AND day = $2;

-- name: CreateTrainJourney :one
INSERT INTO train_journey (train_id , journey_date ,schedule_id, status )
VALUES( $1 , $2 , $3 , $4) RETURNING *;


-- name: GetCoachesByTrain :many
SELECT * FROM coach WHERE trainId = $1;

-- name: GetSeatsByCoach :many
SELECT * FROM seat WHERE coachId = $1;

-- name: GetSeatsByTrain :many
SELECT s.*
FROM seat s 
JOIN coach c ON s.coachId = c.id
WHERE c.trainId = $1;

-- name: GetAllTrain :many
SELECT t.* , ts.*
FROM train t
JOIN train_schedule ts ON t.id = ts.trainid;



-- name: GetAvailableSeats :many
SELECT
    si.coach_type,
    COUNT(*) FILTER (WHERE si.status = 'AVAILABLE') AS available_seats,
    COUNT(*) FILTER (WHERE si.status = 'CONFIRMED') AS booked_seats,
    COUNT(*) AS total_seats
FROM seat_inventory si
JOIN train_journey tj ON tj.id = si.journey_id
WHERE tj.train_id = sqlc.arg(train_id)
  AND tj.journey_date = sqlc.arg(journey_date)
  AND si.quota = sqlc.arg(quota)
GROUP BY si.coach_type
ORDER BY si.coach_type;


-- SELECT *
-- FROM get_available_seats(1, '2026-01-15');

-- name: ValidateTrain :one
SELECT COUNT(*)
FROM train WHERE id = $1;


-- name: ValidateSchedule :one
SELECT COUNT(*)
FROM train_schedule
WHERE trainId = $1
AND day = $2;
 
-- name: ValidateSeatsBelongToTrain :one
SELECT COUNT(*) = $1::int as all_seat_belong,
 count(*)::int as seat_found,
 $1 - COUNT(*) as seats_not_found
 FROM seat s 
 JOIN coach c on s.coachId = c.id
 WHERE s.id = ANY($2::int[])
  AND c.trainId = $3;
        

-- below are not applied till now

-- name: HoldSeat :exec
UPDATE seat_inventory
SET status = 'HELD',
    booking_id = $3
WHERE journey_id = $1
  AND seat_id = $2
  AND status = 'AVAILABLE';

-- name: ConfirmSeat :exec
UPDATE seat_inventory
SET status = 'CONFIRMED'
WHERE booking_id = $1
AND status = 'HELD';


-- name: ReleaseExpiredSeats :exec
UPDATE seat_inventory
SET status = 'AVAILABLE',
    booking_id = NULL
WHERE status = 'HELD'
AND booking_id IN (
    SELECT id FROM booking
    WHERE status = 'PENDING'
    AND createdAt < now() - interval '5 minutes'
);

-- name: GetNextCoachNumber :one
SELECT COALESCE(MAX(coachNumber), 0) + 1
FROM coach
WHERE trainId = $1
FOR UPDATE;

-- name: LockTrainForLayout :one
SELECT id
FROM train
WHERE id = $1
FOR UPDATE;


-- name: GetTrainJourneyById :one
select *
from train_journey
where id = $1;

-- name: LockAvailableSeats :many
SELECT seat_id
FROM seat_inventory
WHERE journey_id = sqlc.arg(journey_id)
  AND coach_type = sqlc.arg(coach_type)
  AND quota = sqlc.arg(quota)
  AND status = 'AVAILABLE'
ORDER BY seat_id
FOR UPDATE SKIP LOCKED
LIMIT sqlc.arg(seat_limit);
