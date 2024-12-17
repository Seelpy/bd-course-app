-- +goose Up
-- +goose StatementBegin


CREATE TABLE author
(
    author_id   BINARY(16) NOT NULL,   -- UUID автора
    avatar_id   BINARY(16),            -- UUID аватара автора
    first_name  VARCHAR(255) NOT NULL, -- Имя автора
    second_name VARCHAR(255) NOT NULL, -- Фамилия автора
    middle_name VARCHAR(255),          -- Отчество автора
    nickname    VARCHAR(255),          -- Никнейм автора
    PRIMARY KEY (author_id),           -- Первичный ключ
    CONSTRAINT fk_author_avatar FOREIGN KEY (avatar_id) REFERENCES image (image_id)
)
    ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


DROP TABLE author;
-- +goose StatementEnd
