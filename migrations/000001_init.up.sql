CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY,
	email VARCHAR(255) NOT NULL UNIQUE,
	password_hash VARCHAR(255) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rooms (
	id UUID PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	description TEXT,
	capacity INT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS schedules (
	id UUID PRIMARY KEY,
	room_id UUID NOT NULL,
	start_time TIME NOT NULL,
	end_time TIME NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	UNIQUE(room_id),
	FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS schedule_days (
	schedule_id UUID NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
	day_of_week INT NOT NULL CHECK (day_of_week > 0 AND day_of_week <= 7),

	PRIMARY KEY (schedule_id, day_of_week)
);

CREATE TABLE IF NOT EXISTS slots (
	id UUID PRIMARY KEY,
	room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
	start_at TIMESTAMPTZ NOT NULL,
	end_at TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	UNIQUE(room_id, start_at, end_at),
	CHECK(start_at < end_at)
);

CREATE TABLE IF NOT EXISTS bookings (
	id UUID PRIMARY KEY,
	slot_id UUID NOT NULL REFERENCES slots(id) ON DELETE CASCADE,
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	status VARCHAR(50) NOT NULL,
	conference_link TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
	CHECK (status IN ('active', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);

CREATE INDEX IF NOT EXISTS idx_slots_room_id ON slots(room_id, start_at);
