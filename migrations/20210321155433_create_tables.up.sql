
CREATE TABLE transactions(
    id serial PRIMARY KEY,
    amount DECIMAL(1000, 10),
    date TIMESTAMP,
    created_at timestamp
);

CREATE TABLE history(
    id SERIAL PRIMARY KEY,
    hour TIMESTAMP,
    amount DECIMAL(1000, 10),
    created_at timestamp,
    updated_at timestamp
);

CREATE INDEX hour_idx ON history (hour);
