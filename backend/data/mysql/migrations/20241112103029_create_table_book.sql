-- +goose Up
-- +goose StatementBegin


CREATE TABLE book
(
    book_id     BINARY(16) NOT NULL,                 -- UUID книги
    cover_id    BINARY(16) DEFAULT NULL,             -- UUID обложки (может быть NULL)
    description TEXT,                                -- Описание книги
    title       VARCHAR(255) NOT NULL,               -- Заголовок книги
    is_publish  BOOLEAN      NOT NULL DEFAULT FALSE, -- Флаг опубликования книги
    PRIMARY KEY (book_id)                           -- Первичный ключ
) ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


DROP TABLE book;
-- +goose StatementEnd
