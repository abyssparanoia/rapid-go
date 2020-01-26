\c catharsis_authentication;

create extension "uuid-ossp";

create table users (
  id uuid DEFAULT uuid_generate_v4 () primary key,
  password varchar NOT NULL,
  display_name varchar NOT NULL,
  icon_image_path varchar NOT NULL,
  background_image_path varchar NOT NULL,
  profile varchar,
  email varchar,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL,
  deleted_at timestamp
);