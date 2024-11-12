-- +goose Up
-- +goose StatementBegin


CREATE TABLE book_chapter_translation
(
    book_chapter_id BINARY(16) NOT NULL,          -- UUID главы книги
    translator_id   BINARY(16) NOT NULL,          -- UUID переводчика
    text            TEXT NOT NULL,                -- Текст перевода
    PRIMARY KEY (book_chapter_id, translator_id), -- Композитный первичный ключ
    CONSTRAINT fk_book_chapter_translation FOREIGN KEY (book_chapter_id) REFERENCES book_chapter (book_chapter_id),
    CONSTRAINT fk_translator FOREIGN KEY (translator_id) REFERENCES user (user_id)
)
    ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


DROP TABLE book_chapter_translation;
-- +goose StatementEnd
