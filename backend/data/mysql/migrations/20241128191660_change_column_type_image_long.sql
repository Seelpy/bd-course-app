-- +goose Up
-- +goose StatementBegin
ALTER TABLE image
    MODIFY COLUMN path LONGTEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE image
    MODIFY COLUMN path TEXT NOT NULL;
-- +goose StatementEnd