-- +goose Up
-- +goose StatementBegin
ALTER TABLE image
    MODIFY COLUMN path TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE image
    MODIFY COLUMN path VARCHAR(255) NOT NULL;
-- +goose StatementEnd