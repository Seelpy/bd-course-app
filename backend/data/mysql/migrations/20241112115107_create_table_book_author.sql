-- +goose Up
-- +goose StatementBegin


CREATE TABLE book_author
(
    book_id   BINARY(16),             -- UUID книги
    author_id BINARY(16),             -- UUID автора
    PRIMARY KEY (book_id, author_id), -- Композитный первичный ключ
    CONSTRAINT fk_author_book FOREIGN KEY (book_id) REFERENCES book (book_id) ON DELETE CASCADE,
    CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES author (author_id) ON DELETE CASCADE
)
    ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


DRop TABLE book_author;
-- +goose StatementEnd
