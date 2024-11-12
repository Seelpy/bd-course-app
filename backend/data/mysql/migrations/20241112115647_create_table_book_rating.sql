-- +goose Up
-- +goose StatementBegin


CREATE TABLE book_rating
(
    book_id BINARY(16) NOT NULL,    -- UUID книги
    user_id BINARY(16) NOT NULL,    -- UUID пользователя
    value   INT NOT NULL,           -- Оценка
    PRIMARY KEY (book_id, user_id), -- Композитный первичный ключ,
    CONSTRAINT fk_rating_book FOREIGN KEY (book_id) REFERENCES book (book_id),
    CONSTRAINT fk_rating_user FOREIGN KEY (user_id) REFERENCES user (user_id)
)
    ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


DROP TABLE book_rating;
-- +goose StatementEnd
