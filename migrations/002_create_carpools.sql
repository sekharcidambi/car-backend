-- Create table for carpools
CREATE TABLE carpools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_id VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    schedule_date TIMESTAMP,
    schedule_time TIMESTAMP,
    recurring_option VARCHAR(50),
    available_seats INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create table for carpool stops
CREATE TABLE carpool_stops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    carpool_id UUID REFERENCES carpools(id) ON DELETE CASCADE,
    address TEXT NOT NULL,
    stop_order INTEGER NOT NULL,  -- 0 for starting point, highest number for destination
    stop_type VARCHAR(20) CHECK (stop_type IN ('START', 'INTERMEDIATE', 'DESTINATION')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster lookups
CREATE INDEX idx_carpool_stops_carpool_id ON carpool_stops(carpool_id);
CREATE INDEX idx_carpool_stops_order ON carpool_stops(carpool_id, stop_order);

-- Optional: Create an update trigger for updated_at (carpools table)
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

-- Add trigger for carpool_stops table
CREATE TRIGGER update_carpool_stops_updated_at
    BEFORE UPDATE ON carpool_stops
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();