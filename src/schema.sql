CREATE TABLE urls (
    id SERIAL PRIMARY KEY, 
    long_url TEXT NOT NULL, 
    short_url VARCHAR(8) NOT NULL UNIQUE DEFAULT substr(md5(random()::text), 0, 8), 
    date_created TIME NOT NULL DEFAULT now(),
    hit_count INT NOT NULL DEFAULT 0,
    last_hit TIMESTAMP,
    user_agent TEXT
);