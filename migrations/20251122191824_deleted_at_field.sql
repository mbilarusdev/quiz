-- +goose Up
-- +goose StatementBegin
ALTER TABLE questions ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE answers ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE questions DROP COLUMN deleted_at;
ALTER TABLE answers DROP COLUMN deleted_at;
-- +goose StatementEnd
