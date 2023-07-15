CREATE TABLE if not EXISTS football_players (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255),
    last_name VARCHAR(255),
    value INTEGER,
    team VARCHAR(255)
);