# sensors_backend

Use the following script in your local POSTGRE client 🤖

```-- Database: density

-- DROP DATABASE IF EXISTS density;

CREATE DATABASE IF NOT EXISTS density
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'C'
    LC_CTYPE = 'C'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;


-- SCHEMA: sensor

-- DROP SCHEMA IF EXISTS sensor ;

CREATE SCHEMA IF NOT EXISTS sensor
    AUTHORIZATION postgres;

— TABLE europe_sensor

CREATE TABLE sensor.europe_sensors(
	id serial primary key,
	sensor_id VARCHAR(255) NOT NULL,
	room_id VARCHAR(255) NOT NULL,
	floor_id VARCHAR(255) NOT NULL,
	building_id VARCHAR(255) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP)
  ```
The repo with sensor is here https://github.com/davidfradel/sensors 💓
