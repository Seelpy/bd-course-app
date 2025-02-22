CREATE TABLE book_genre
(
    book_id  BINARY(16),           -- UUID книги
    genre_id BINARY(16),           -- UUID жанра
    PRIMARY KEY (book_id, genre_id), -- Композитный первичный ключ  
    CONSTRAINT fk_genre_book FOREIGN KEY (book_id) REFERENCES book (book_id),
    CONSTRAINT fk_genre FOREIGN KEY (genre_id) REFERENCES genre (genre_id)
) ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;