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

-- name: GetAvailableSeats :many
CREATE OR REPLACE FUNCTION get_avaliable_seats(
    p_train_id INTEGER,
    p_travel_day TIME
)
RETURNS TABLE (
    coach_type TEXT,
    total_seats BIGINT,
    booked_seats BIGINT,
    available_seats BIGINT
)
LANGUAGE plpgsql
AS $$ 
BEGIN 
    RETURN QUERY 
    SELECT c.coach_type ,
           count(s.id) ::BIGINT as total_seats
           count(bi.seatId) :: BIGINT as booked_seats
          (count(s.id)- count(bi.seatId)) :: BIGINT as available_seats
    FROM train t 
    JOIN coach c on t.id = c.trainId
    JOIN seat s on  c.id = s.coachId
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
        


