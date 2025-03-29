-- +goose Up
-- +goose StatementBegin
ALTER TABLE available_characters ADD COLUMN level integer NOT NULL DEFAULT 1;
ALTER TABLE available_characters ADD COLUMN level_experience integer NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE available_characters DROP COLUMN level;
ALTER TABLE available_characters DROP COLUMN level_experience;
-- +goose StatementEnd
