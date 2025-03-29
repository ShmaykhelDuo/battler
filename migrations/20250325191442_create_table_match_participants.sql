-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS match_participants (
    user_id uuid REFERENCES users (id),
    match_id uuid REFERENCES matches (id),
    character_number integer NOT NULL,
    result integer NOT NULL,
    PRIMARY KEY (user_id, match_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS match_participants;
-- +goose StatementEnd
