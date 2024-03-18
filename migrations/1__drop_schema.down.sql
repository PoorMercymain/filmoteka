BEGIN;
DROP TABLE IF EXISTS actors;

DROP TABLE IF EXISTS films;

DROP TABLE IF EXISTS film_actor;

DROP TABLE IF EXISTS auth;

DROP TRIGGER IF EXISTS check_actor_birthday_before_insert ON film_actor;
DROP FUNCTION IF EXISTS check_actor_birthday_before_film_release();
COMMIT;