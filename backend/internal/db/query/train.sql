-- name: CreateTrain :one
INSERT into train(trainNumber,trainName,source,destination )
VALUES ( $1 ,$2 ,$3,$4) 
RETURNING *;

-- name: CreateTrainSchedule :one
INSERT into trainSchedule(trainId,day,arrivalTime,departureTime)
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
SELECT * FROM trainSchedule
WHERE trainId = $1 AND day = $2;

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
JOIN trainSchedule ts ON t.id = ts.trainid;



-- name: GetAvailableSeats :many
CREATE  OR REPLACE FUNCTION get_available_seats(
    p_train_id INTEGER,
    p_travel_date DATE  -- Changed from TIME to DATE to match your schema
)
RETURNS TABLE (
    coach_type coach_type,
    total_seats BIGINT,
    booked_seats BIGINT,
    available_seats BIGINT
)
LANGUAGE plpgsql
AS $$ 
BEGIN 
    RETURN QUERY 
    SELECT c.coachtype,
           COUNT(s.id)::BIGINT as total_seats,
           COUNT(bi.seatId)::BIGINT as booked_seats,
           (COUNT(s.id) - COUNT(bi.seatId))::BIGINT as available_seats
    FROM train t 
    JOIN coach c ON t.id = c.trainId
    JOIN seat s ON c.id = s.coachId
    LEFT JOIN (
        SELECT DISTINCT bi.seatId
        FROM booking b
        JOIN bookingItem bi ON b.id = bi.bookingId
        WHERE b.trainId = p_train_id
          AND b.travelDate = p_travel_date
          AND b.status IN ('CONFIRMED', 'PENDING')
    ) bi ON s.id = bi.seatId
    WHERE t.id = p_train_id
    GROUP BY c.coachtype
    ORDER BY c.coachtype;
END;
$$;

-- name: GetAvailableSeatsExecute :many
SELECT * FROM get_available_seats($1, $2);

-- name: ValidateTrain :one
SELECT COUNT(*)
FROM train WHERE id = $1;

-- name: ValidateSchedule :one
SELECT count(*)
FROM trainSchedule ts
WHERE ts.trainId = $1 
AND DATE(ts.arrivaltime) = $2::DATE;

-- name: ValidateSeatsBelongToTrain :one
SELECT COUNT(*) = $1::int as all_seat_belong,
 count(*)::int as seat_found,
 $1 - COUNT(*) as seats_not_found
 FROM seat s 
 JOIN coach c on s.coachId = c.id
 WHERE s.id = ANY($2::int[])
  AND c.trainId = $3;
        


