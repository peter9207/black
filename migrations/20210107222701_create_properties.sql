-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS properties (
  id UUID PRIMARY KEY,
  type varchar,
  value numeric,
  name varchar,
  meta varchar,
  created_at timestamp
  )

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
  drop table properties
