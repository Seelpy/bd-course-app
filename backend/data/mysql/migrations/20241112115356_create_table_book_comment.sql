-- +goose Up
-- +goose StatementBegin


CREATE TABLE book_comment
(
    book_comment_id BINARY(16) NOT NULL,                -- UUID комментария к книге
    book_id         BINARY(16) NOT NULL,                -- UUID книги
    user_id         BINARY(16) NOT NULL,                -- UUID пользователя
    comment         TEXT NOT NULL,                      -- Текст комментария
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP, -- Дата и время создания комментария
    PRIMARY KEY (book_comment_id),                      -- Первичный ключ,
    CONSTRAINT fk_comment_book FOREIGN KEY (book_id) REFERENCES book (book_id),
    CONSTRAINT fk_comment_user FOREIGN KEY (user_id) REFERENCES user (user_id)
)
    ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


DROP TABLE book_comment;
-- +goose StatementEnd
