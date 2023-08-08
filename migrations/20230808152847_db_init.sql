-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS USERS(
    ID SERIAL PRIMARY KEY,
    LOGIN TEXT NOT NULL,
    PASSWORD TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE USERS;
-- +goose StatementEnd
