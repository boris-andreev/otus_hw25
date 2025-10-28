-- +goose Up
-- +goose StatementBegin
ALTER TABLE homeworks ADD COLUMN description TEXT;
ALTER TABLE studies ADD COLUMN topic TEXT;
ALTER TABLE workouts ADD COLUMN target TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE homeworks DROP COLUMN description;
ALTER TABLE studies DROP COLUMN topic;
ALTER TABLE workouts DROP COLUMN target;
-- +goose StatementEnd
