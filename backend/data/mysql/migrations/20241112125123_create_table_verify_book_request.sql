-- +goose Up
-- +goose StatementBegin
CREATE TABLE verify_book_request
(
    verify_book_request_id BINARY(16) NOT NULL,                -- UUID запроса на проверку книги
    translator_id          BINARY(16),                         -- UUID переводчика
    book_id                BINARY(16),                         -- UUID книги
    is_verified            BOOLEAN  DEFAULT NULL,              -- Флаг проверки
    send_date              DATETIME DEFAULT CURRENT_TIMESTAMP, -- Дата отправки запроса
    PRIMARY KEY (verify_book_request_id),                      -- Первичный ключ
    CONSTRAINT fk_verify_translator FOREIGN KEY (translator_id) REFERENCES user (user_id),
    CONSTRAINT fk_verify_book FOREIGN KEY (book_id) REFERENCES book (book_id)
)
    ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE verify_book_request;
-- +goose StatementEnd
