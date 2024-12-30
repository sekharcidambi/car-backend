-- Create table users
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255),
    city VARCHAR(255),
    state VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
); 

-- Create table user analytics
CREATE TABLE user_analytics (
    user_id UUID PRIMARY KEY,
    number_of_carpools INTEGER DEFAULT 0,
    number_of_completed_rides INTEGER DEFAULT 0,
    number_of_completed_rides_as_driver INTEGER DEFAULT 0,
    driving_miles_saved INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

ALTER TABLE users 
ADD COLUMN clerk_id TEXT NOT NULL,
ADD CONSTRAINT unique_clerk_id UNIQUE (clerk_id);