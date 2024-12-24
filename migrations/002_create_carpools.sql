-- Create table for carpools
CREATE TABLE carpools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_id VARCHAR(255),
    schedule_date TIMESTAMP,
    schedule_time TIMESTAMP,
    recurring_option VARCHAR(50),
    starting_point_lat FLOAT,
    starting_point_lng FLOAT,
    destination_lat FLOAT,
    destination_lng FLOAT,
    available_seats INTEGER,
    music_preference VARCHAR(50),
    smoking_allowed BOOLEAN,
    pets_allowed BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Optional: Create an update trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_carpools_updated_at
    BEFORE UPDATE ON carpools
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();