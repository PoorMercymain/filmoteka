BEGIN;
CREATE INDEX IF NOT EXISTS actors_idx ON actors USING BTREE(id, name, gender, birthday);
CREATE INDEX IF NOT EXISTS films_idx ON films USING BTREE(id, title, description, release_date, rating);
COMMIT;