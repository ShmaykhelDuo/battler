-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS currency_conversions (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users (id),
    started_at timestamptz NOT NULL,
    finishes_at timestamptz NOT NULL,
    target_currency_id integer NOT NULL,
    amount bigint NOT NULL,
    is_claimed boolean NOT NULL DEFAULT false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS currency_conversions;
-- +goose StatementEnd
