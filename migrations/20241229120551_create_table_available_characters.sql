-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS available_characters (
    user_id uuid REFERENCES users (id),
    number integer,
    PRIMARY KEY (user_id, number)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS available_characters;
-- +goose StatementEnd
