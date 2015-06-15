DROP SCHEMA tinyplannr_api CASCADE;
CREATE SCHEMA tinyplannr_api;

DROP TABLE IF EXISTS tinyplannr_api.user_api;
CREATE TABLE tinyplannr_api.user_api (
  user_id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  zip_code INTEGER,
  is_active BOOLEAN,
  create_dt TIMESTAMP,
  update_dt TIMESTAMP
);

DROP TABLE IF EXISTS tinyplannr_api.user_auth;
CREATE TABLE tinyplannr_api.user_auth (
  auth_id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES tinyplannr_api.user_api (user_id),
  email VARCHAR(255) REFERENCES tinyplannr_api.user_api (email) UNIQUE,
  hash_pw TEXT,
  created_dt TIMESTAMP,
  update_dt TIMESTAMP,
  last_login_dt TIMESTAMP
);

DROP TABLE IF EXISTS tinyplannr_api.event;
CREATE TABLE tinyplannr_api.event (
  event_id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES tinyplannr_api.user_api (user_id),
  title TEXT,
  description TEXT,
  location TEXT,
  all_day BOOLEAN,
  start_dt TIMESTAMP,
  end_dt TIMESTAMP,
  create_dt TIMESTAMP,
  update_dt TIMESTAMP
);