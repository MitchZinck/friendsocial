CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    zip_code VARCHAR(20),
    country VARCHAR(100) NOT NULL,
    latitude DECIMAL(9, 6),
    longitude DECIMAL(9, 6)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    location_id INTEGER,
    CONSTRAINT uq_email UNIQUE (email),
    CONSTRAINT fk_location FOREIGN KEY (location_id) REFERENCES locations (id)
);

CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_location ON users (location_id);

CREATE TABLE activities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    estimated_time INTERVAL NOT NULL,
    location_id INTEGER NOT NULL,
    user_created BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_location_id FOREIGN KEY (location_id)
    REFERENCES locations (id)
);

CREATE TABLE scheduled_activities (
    id SERIAL PRIMARY KEY,
    activity_id INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    scheduled_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_activity_id FOREIGN KEY (activity_id)
    REFERENCES activities (id)
);

CREATE INDEX idx_scheduled_activities_activity_id ON scheduled_activities (activity_id); -- Index on activity_id

CREATE TABLE friends (
    user_id INTEGER NOT NULL,
    friend_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_ordered_id1 INTEGER GENERATED ALWAYS AS (LEAST(user_id, friend_id)) STORED,
    user_ordered_id2 INTEGER GENERATED ALWAYS AS (GREATEST(user_id, friend_id)) STORED,
    CONSTRAINT pk_friends PRIMARY KEY (user_id, friend_id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_friend FOREIGN KEY (friend_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT chk_not_self_friend CHECK (user_id <> friend_id),
    CONSTRAINT uq_friends_pair UNIQUE (user_ordered_id1, user_ordered_id2) -- Prevent duplicate relationships
);

CREATE INDEX idx_friends_user_id ON friends (user_id); -- Index on user_id
CREATE INDEX idx_friends_friend_id ON friends (friend_id); -- Index on friend_id
CREATE INDEX idx_friends_pair ON friends (user_ordered_id1, user_ordered_id2); -- Index on pair

CREATE TABLE user_availability (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    day_of_week VARCHAR(9) NOT NULL, -- e.g., 'Monday', 'Tuesday'
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    is_available BOOLEAN DEFAULT false,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
    REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_user_availability_user_id ON user_availability (user_id); -- Index on user_id

CREATE TABLE user_activity_preferences (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    activity_id INTEGER NOT NULL,
    frequency INTEGER NOT NULL,
    frequency_period VARCHAR(50) NOT NULL, -- e.g., 'daily', 'weekly', 'monthly'
    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
    REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_activity_id FOREIGN KEY (activity_id)
    REFERENCES activities (id) ON DELETE CASCADE
);

CREATE INDEX idx_user_activity_preferences_user_id ON user_activity_preferences (user_id); -- Index on user_id
CREATE INDEX idx_user_activity_preferences_activity_id ON user_activity_preferences (activity_id); -- Index on activity_id

CREATE TABLE activity_participants (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    scheduled_activity_id INTEGER NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
    REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_scheduled_activity_id FOREIGN KEY (scheduled_activity_id)
    REFERENCES scheduled_activities (id) ON DELETE CASCADE,
    CONSTRAINT uq_activity_user UNIQUE (user_id, scheduled_activity_id)
);

CREATE INDEX idx_activity_participants_user_id ON activity_participants (user_id);
CREATE INDEX idx_activity_participants_scheduled_activity_id ON activity_participants (scheduled_activity_id);
