-- +goose Up
-- +goose StatementBegin


CREATE TABLE book_chapter
(
    book_chapter_id BINARY(16) NOT NULL,   -- UUID главы книги
    book_id         BINARY(16) NOT NULL,   -- UUID книги
    chapter_index   INT UNIQUE   NOT NULL, -- Индекс главы (уникальный)
    title           VARCHAR(255) NOT NULL, -- Заголовок главы
    PRIMARY KEY (book_chapter_id),         -- Первичный ключ
    CONSTRAINT fk_book_chapter FOREIGN KEY (book_id) REFERENCES book (book_id) ON DELETE CASCADE
)
    ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


DROP TABLE book_chapter;
-- +goose StatementEnd
