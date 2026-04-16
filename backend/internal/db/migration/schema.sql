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
    'WAITLIST',
    'CANCELLED',
    'EXPIRED'
);

CREATE TYPE booking_type as ENUM (
    'NORMAL',
    'WAITLIST',
    'TATKAL'
);

CREATE type waiting_status as ENUM (
    'WAITING',
    'CONFIRMED',
    'CANCELLED'
);

CREATE TYPE payment_status AS ENUM (
    'PENDING',
    'FAILED',
    'SUCCESS'
);

CREATE TYPE refund_status AS ENUM (
    'PENDING',
    'FAILED',
    'SUCCESS'
);

CREATE TYPE seat_status AS ENUM (
  'AVAILABLE',
  'HELD',
  'CONFIRMED'
);

CREATE TYPE journey_status AS ENUM ('OPEN','CHARTED','CANCELLED');


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
 -- things are changingn here
CREATE TABLE train_schedule (
    id  SERIAL PRIMARY KEY,
    trainId INTEGER REFERENCES train(id) ON DELETE CASCADE,
    day day_of_week NOT NULL,
    arrivalTime  TIMESTAMPTZ NOT NULL,
    departureTime TIMESTAMPTZ NOT NULL,
    CONSTRAINT unique_train_schedule UNIQUE(trainId,day)
);

-- ALTER TABLE train_schedule 
--     ALTER COLUMN arrivalTime TYPE TIME,
--     ALTER COLUMN arrivalTime SET NOT NULL,
--     ALTER COLUMN departureTime TYPE TIME,
--     ALTER COLUMN departureTime SET NOT NULL;

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

CREATE TABLE train_journey (
    id SERIAL PRIMARY KEY,
    train_id INT REFERENCES train(id),
    journey_date DATE NOT NULL,
    schedule_id INT REFERENCES train_schedule(id),
    status journey_status DEFAULT 'OPEN',
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE(train_id, journey_date)
);




CREATE TYPE seat_quota AS ENUM ('NORMAL', 'TATKAL');

-- made for working on tatkal not implemented yet
CREATE TABLE seat_inventory (
  journey_id INT REFERENCES train_journey(id),
  seat_id INT REFERENCES seat(id),
  coach_type coach_type NOT NULL,
  quota seat_quota NOT NULL,
  status seat_status NOT NULL DEFAULT 'AVAILABLE',
  booking_id INT REFERENCES booking(id) ON DELETE SET NULL,
  PRIMARY KEY (journey_id, seat_id)
);

-- i have changed the tatkal schema

CREATE Table tatkal_config (
    id SERIAL PRIMARY KEY,
    train_id INTEGER REFERENCES train(id) ON Delete CASCADE,
    coach_type coach_type NOT NULL,
    tatkal_start_time TIMESTAMP NOT NULL,
    tatkal_end_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE(train_id, coach_type)
);

CREATE TABLE tatkal_waitlist (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) not null,
    coach_type coach_type NOT NULL,
    journey_id INT REFERENCES train_journey(id),
    wl_position INT NOT NULL,
    status waiting_status DEFAULT 'WAITING',
    created_at TIMESTAMP DEFAULT now(),
   UNIQUE(user_id, journey_id)
);




CREATE TABLE payment (
    id SERIAL PRIMARY KEY,
    bookingId INT REFERENCES booking(id) ON DELETE RESTRICT,
    amount FLOAT NOT NULL,
    status payment_status DEFAULT 'PENDING',
    transactionId TEXT NOT NULL,
    createdAt TIMESTAMP DEFAULT now()
);

CREATE TABLE booking (
    id SERIAL PRIMARY KEY,
    userId UUID REFERENCES users(id) ON DELETE RESTRICT,
    journey_id INT REFERENCES train_journey(id) ON DELETE RESTRICT,
    booking_type booking_type NOT NULL DEFAULT 'NORMAL',
    status booking_status NOT NULL DEFAULT 'PENDING',
    holdToken TEXT UNIQUE,
    createdAt TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE bookingItem (
    id SERIAL PRIMARY KEY,
    bookingId INT REFERENCES booking(id) ON DELETE CASCADE,
    seatId INT REFERENCES seat(id) ON DELETE RESTRICT,
    bookingStatus booking_status NOT NULL DEFAULT 'PENDING',
    UNIQUE (bookingId, seatId)
);


-- This enforces that the payment must always reference a valid booking.
ALTER TABLE payment
ADD CONSTRAINT fk_payment_booking
FOREIGN KEY (bookingId) 
REFERENCES booking(id)
ON DELETE RESTRICT;


-- not made table till now
CREATE TABLE booking_passenger (
    id SERIAL PRIMARY KEY,
    booking_id INT REFERENCES booking(id) ON DELETE CASCADE,
    seat_id INT REFERENCES seat(id) ON DELETE RESTRICT,
    name TEXT NOT NULL,
    age INT NOT NULL,
    gender TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE (booking_id, seat_id)
);


CREATE TABLE Refund (
    id SERIAL PRIMARY KEY,
    userId uuid REFERENCES users(id) on delete CASCADE,
    bookingId INTEGER REFERENCES booking(id) on delete CASCADE,
    amount INTEGER NOT NULL,
    status refund_status not null DEFAULT 'PENDING',
    createdAt TIMESTAMP NOT NULL DEFAULT now(),
    updatedAt TIMESTAMP NOT NULL DEFAULT now()

);

CREATE TABLE waitlist (
    id SERIAL PRIMARY KEY,
    journey_id INT REFERENCES train_journey(id),
    bookingId INT REFERENCES booking(id),
    waitlist_number INTEGER not NULL,
    status  waiting_status NOT NULL DEFAULT 'WAITING',
    priority_level INTEGER DEFAULT 10,
    createdAt TIMESTAMP not NULL DEFAULT now(),
    updatedAt TIMESTAMP not NULL DEFAULT now()
);

-- should me moved to top




-- indexex 
CREATE UNIQUE INDEX idx_train_number ON train(trainNumber);
CREATE INDEX idx_train_src ON train(source);
CREATE INDEX idx_train_dest ON train(destination);

CREATE INDEX idx_booking_user ON booking(userId);
CREATE INDEX idx_booking_status ON booking(status);
CREATE INDEX idx_payment_status ON payment(status);
CREATE INDEX idx_booking_journey ON booking(journey_id);
CREATE INDEX idx_inventory_search
ON seat_inventory (journey_id, coach_type, quota, status);

CREATE INDEX idx_inventory_booking
ON seat_inventory (booking_id);


