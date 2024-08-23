CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    CONSTRAINT uq_email UNIQUE (email) -- Unique constraint on email
);

CREATE INDEX idx_users_email ON users (email); -- Index on email

CREATE TABLE activity_locations (
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

CREATE TABLE activities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    estimated_time INTERVAL NOT NULL,
    location_id INTEGER NOT NULL,
    CONSTRAINT fk_location_id FOREIGN KEY (location_id)
    REFERENCES activity_locations (id)
);

CREATE TABLE user_activities (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    activity_id INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
    REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_activity_id FOREIGN KEY (activity_id)
    REFERENCES activities (id)
);

CREATE INDEX idx_user_activities_user_id ON user_activities (user_id); -- Index on user_id
CREATE INDEX idx_user_activities_activity_id ON user_activities (activity_id); -- Index on activity_id

CREATE TABLE friends (
    user_id INTEGER NOT NULL,
    friend_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_friends PRIMARY KEY (user_id, friend_id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_friend FOREIGN KEY (friend_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT chk_not_self_friend CHECK (user_id <> friend_id),
    CONSTRAINT uq_friends_pair UNIQUE (LEAST(user_id, friend_id), GREATEST(user_id, friend_id)) -- Prevent duplicate relationships
);

CREATE INDEX idx_friends_user_id ON friends (user_id); -- Index on user_id
CREATE INDEX idx_friends_friend_id ON friends (friend_id); -- Index on friend_id

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

CREATE TABLE manual_activities (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    activity_id INTEGER NULL, -- nullable to allow for truly custom activities
    name VARCHAR(100) NOT NULL,
    description TEXT NULL, -- nullable if description is optional
    estimated_time INTERVAL NULL, -- nullable to allow flexible time entries
    location_id INTEGER NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_location_id FOREIGN KEY (location_id) REFERENCES activity_locations (id) ON DELETE SET NULL,
    CONSTRAINT fk_activity_id FOREIGN KEY (activity_id) REFERENCES activities (id)
);

CREATE INDEX idx_manual_activities_user_id ON manual_activities (user_id); -- Index on user_id

-- New table for participants, which can reference either manual or automated activities
CREATE TABLE activity_participants (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    activity_id INTEGER NULL,
    manual_activity_id INTEGER NULL,
    is_creator BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
    REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_activity_id FOREIGN KEY (activity_id)
    REFERENCES activities (id) ON DELETE CASCADE,
    CONSTRAINT fk_manual_activity_id FOREIGN KEY (manual_activity_id)
    REFERENCES manual_activities (id) ON DELETE CASCADE,
    CONSTRAINT chk_activity_participant CHECK (
        (activity_id IS NOT NULL AND manual_activity_id IS NULL)
        OR (manual_activity_id IS NOT NULL AND activity_id IS NULL)
    ), -- Ensure that only one of activity_id or manual_activity_id is filled
    CONSTRAINT uq_activity_user UNIQUE (user_id, activity_id, manual_activity_id)
);

CREATE INDEX idx_activity_participants_user_id ON activity_participants (user_id); -- Index on user_id
CREATE INDEX idx_activity_participants_activity_id ON activity_participants (activity_id); -- Index on activity_id
CREATE INDEX idx_activity_participants_manual_activity_id ON activity_participants (manual_activity_id); -- Index on manual_activity_id
