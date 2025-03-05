-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users_friends (
    user_id uuid REFERENCES users (id),
    friend_id uuid REFERENCES users (id),
    PRIMARY KEY (user_id, friend_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users_friends;
-- +goose StatementEnd
