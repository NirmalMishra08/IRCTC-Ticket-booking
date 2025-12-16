CREATE TYPE user_role AS ENUM(
    'ADMIN',
    'USER'
);

CREATE TYPE provider as ENUM (
    'EMAIL',
    'APPLE',
    'PASSWORD'
);




CREATE TYPE day_of_week AS ENUM ('MON', 'TUE', 'WED', 'THU', 'FRI', 'SAT','SUN');

CREATE type coach_type as ENUM ('3A','2A','1A','SL','GN');

CREATE type berth_type as ENUM ('UP','DOWN','MID');

CREATE TYPE booking_status AS ENUM (
    'PENDING',
    'CONFIRMED',
    'CANCELLED',
    'EXPIRED'
);

CREATE TYPE payment_status AS ENUM (
    'PENDING',
    'FAILED',
    'SUCCESS'
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fullname TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    role user_role NOT NULL DEFAULT 'USER',
    password_hash TEXT ,
    provider provider NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

ALTER table users ADD COLUMN phone TEXT ;


CREATE TABLE train (
    id SERIAL PRIMARY KEY,
    trainNumber INTEGER NOT NULL,
    trainName TEXT NOT NULL,
    source TEXT NOT NULL,
    destination TEXT NOT NULL,
    CONSTRAINT unique_train_service UNIQUE (trainNumber)
);

CREATE TABLE trainSchedule (
    id  SERIAL PRIMARY KEY,
    trainId INTEGER REFERENCES train(id) ON DELETE CASCADE,
    day day_of_week NOT NULL,
    arrivalTime  TIMESTAMPTZ NOT NULL,
    departureTime TIMESTAMPTZ NOT NULL,
    CONSTRAINT unique_train_schedule UNIQUE(trainId,day)
);

CREATE Table tatkal (
    id SERIAL PRIMARY KEY,
    trainId INTEGER REFERENCES train(id) ON Delete CASCADE,
    coachType coach_type NOT NULL,
    totalSeats INT
);

CREATE TABLE coach (
   id SERIAL PRIMARY KEY,
   trainId INTEGER REFERENCES train(id) on delete CASCADE,
   coachtype coach_type NOT NULL,
   coachNumber INTEGER NOT NULL,
   CONSTRAINT unique_coach_per_train UNIQUE (trainId, coachNumber)
);


CREATE TABLE seat (
    id SERIAL PRIMARY KEY,
    coachId  INTEGER REFERENCES coach(id) ON DELETE CASCADE,
    seatno INTEGER NOT NULL,
    berth  berth_type NOT NULL
);

CREATE Table payment (
    id  SERIAL PRIMARY KEY ,
    bookingId  INTEGER,
    amount   FLOAT NOT NULL,
    status payment_status DEFAULT 'PENDING',
    transactionId TEXT NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE booking (
    id SERIAL PRIMARY KEY,
    userId UUID REFERENCES users(id) ON DELETE RESTRICT,
    trainId INTEGER REFERENCES train(id) ON DELETE RESTRICT,
    travelDate DATE NOT NULL,
    status booking_status NOT NULL DEFAULT 'PENDING',
    holdToken TEXT UNIQUE,
    paymentId INTEGER REFERENCES payment(id) ON DELETE RESTRICT,
    createdAt TIMESTAMP NOT NULL DEFAULT now()

);

CREATE TABLE bookingItem (
    id SERIAL PRIMARY KEY,
    bookingId  INTEGER REFERENCES booking(id) ON DELETE CASCADE,
    seatId INT REFERENCES seat(id) ON DELETE RESTRICT,
    trainScheduleId INTEGER REFERENCES trainSchedule(id) ON DELETE RESTRICT,
    CONSTRAINT unique_seat_booked UNIQUE (seatId, trainScheduleId)
);

-- This enforces that the payment must always reference a valid booking.
ALTER TABLE payment
ADD CONSTRAINT fk_payment_booking
FOREIGN KEY (bookingId) 
REFERENCES booking(id) 
ON DELETE RESTRICT;


-- indexex 
CREATE UNIQUE INDEX idx_train_number ON train(trainNumber);
CREATE INDEX idx_train_src ON train(source);
CREATE INDEX idx_train_dest ON train(destination);

CREATE INDEX idx_booking_user ON booking(userId);
CREATE INDEX idx_booking_train_date ON booking(trainId, travelDate);
CREATE INDEX idx_booking_status ON booking(status);
CREATE INDEX idx_payment_status ON payment(status);
CREATE INDEX idx_bookingitem_schedule ON bookingItem(trainScheduleId);
