-- Create table for carpools
CREATE TABLE carpools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_id VARCHAR(255),
    carpool_name VARCHAR(255) NOT NULL, 
    status BOOLEAN DEFAULT TRUE, 
    recurring_option VARCHAR(50),
    available_seats INTEGER,
    destination_address TEXT, 
    seats INTEGER, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create table for carpool invites
CREATE TABLE invites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_user UUID NOT NULL,
    to_user UUID NOT NULL,
    carpool_id UUID NOT NULL,
    message TEXT, 
    status INTEGER, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (from_user) REFERENCES users(id),
    FOREIGN KEY (to_user) REFERENCES users(id),
    FOREIGN KEY (carpool_id) REFERENCES carpools(id)
);

--Create table for carpool members
CREATE TABLE carpool_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    carpool_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (carpool_id) REFERENCES carpools(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE (carpool_id, user_id) -- Ensure unique combination of carpool_id and user_id
);

-- Create table for carpool rides
CREATE TABLE carpool_rides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    carpool_id UUID NOT NULL,
    driver_id UUID NOT NULL, 
    status INTEGER, 
    location_lat FLOAT,
    location_lng FLOAT, 
    miles_saved FLOAT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (carpool_id) REFERENCES carpools(id),
    FOREIGN KEY (driver_id) REFERENCES users(id) 
);


-- Create table for carpool stops(transaction with carpool rides)
CREATE TABLE carpool_stops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    carpool_ride_id UUID NOT NULL,
    address TEXT NOT NULL,
    stop_order INTEGER NOT NULL, 
    stop_type VARCHAR(20) CHECK (stop_type IN ('START', 'INTERMEDIATE', 'DESTINATION')),
    user_id UUID, 
    FOREIGN KEY (carpool_ride_id) REFERENCES carpool_rides(id) ON DELETE CASCADE, 
    FOREIGN KEY (user_id) REFERENCES users(id), 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);





