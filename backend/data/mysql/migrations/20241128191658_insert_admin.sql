-- +goose Up
-- +goose StatementBegin
INSERT INTO `user` (user_id, avatar_id, login, role, password, about_me)
VALUES (UNHEX(REPLACE(UUID(), '-', '')), NULL, 'admin', 0, '12345', 'Administrator account');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM `user` WHERE role = 0;
-- +goose StatementEnd
