-- +migrate Up
-- +Migrate StatementBegin

CREATE TABLE servo (
    id INTERGER PRIMARY KEY,
    status_servo INTERGER
)

-- +Migrate StatementEnd