BEGIN;
CREATE TABLE IF NOT EXISTS actors  (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    gender BOOLEAN,
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

CREATE OR REPLACE FUNCTION check_actor_birthday_before_film_release()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM actors
        WHERE id = NEW.actor_id AND birthday > (
            SELECT release_date FROM films WHERE id = NEW.film_id
        )
    ) THEN
        RAISE EXCEPTION 'Играющие в фильме актеры не могут родиться позже даты выпуска фильма.' USING ERRCODE = 'P0001';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_actor_birthday_before_insert
BEFORE INSERT ON film_actor
FOR EACH ROW EXECUTE FUNCTION check_actor_birthday_before_film_release();
COMMIT;