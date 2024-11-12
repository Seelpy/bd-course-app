CREATE TABLE user_book_favourites
(
    user_id BINARY(16) NOT NULL,    -- UUID пользователя
    book_id BINARY(16) NOT NULL,    -- UUID книги
    type    SMALLINT,               -- Тип избранной книги
    PRIMARY KEY (user_id, book_id), -- Композитный первичный ключ
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES user (user_id),
    CONSTRAINT fk_book FOREIGN KEY (book_id) REFERENCES book (book_id)
) ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;