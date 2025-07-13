CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    film_id INT REFERENCES films(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    rating INT CHECK (rating BETWEEN 1 AND 10) NOT NULL,
    comment TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
); 