-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users_balances (
    user_id uuid REFERENCES users (id),
    currency_id integer,
    amount bigint NOT NULL,
    PRIMARY KEY (user_id, currency_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users_balances;
-- +goose StatementEnd
