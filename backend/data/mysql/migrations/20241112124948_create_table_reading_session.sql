-- +goose Up
-- +goose StatementBegin
CREATE TABLE reading_session
(
    book_id         BINARY(16) NOT NULL,                -- UUID книги
    book_chapter_id BINARY(16) NOT NULL,                -- UUID главы книги
    user_id         BINARY(16) NOT NULL,                -- UUID пользователя
    last_read_time  DATETIME DEFAULT CURRENT_TIMESTAMP, -- Время последнего чтения
    PRIMARY KEY (book_id, book_chapter_id, user_id),    -- Композитный первичный ключ
    CONSTRAINT fk_reading_session_book FOREIGN KEY (book_id) REFERENCES book (book_id) ON DELETE CASCADE,
    CONSTRAINT fk_reading_session_chapter FOREIGN KEY (book_chapter_id) REFERENCES book_chapter (book_chapter_id) ON DELETE CASCADE,
    CONSTRAINT fk_reading_session_user FOREIGN KEY (user_id) REFERENCES user (user_id) ON DELETE CASCADE
)
    ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reading_session;
-- +goose StatementEnd
