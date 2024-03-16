BEGIN;
CREATE TABLE IF NOT EXISTS actors  (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    gender TEXT,
    birthday TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS films (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    release_date TIMESTAMPTZ,
    rating INT
);

CREATE TABLE IF NOT EXISTS film_actor (
    actor_id INT,
    film_id INT,
    PRIMARY KEY (actor_id, film_id),
    FOREIGN KEY (actor_id) REFERENCES actors(id) ON DELETE CASCADE,
    FOREIGN KEY (film_id) REFERENCES films(id) ON DELETE CASCADE
);
COMMIT;