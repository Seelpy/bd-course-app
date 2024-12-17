CREATE TABLE genre
(
    genre_id BINARY(16) NOT NULL,   -- UUID жанра
    name     VARCHAR(255) NOT NULL, -- Название жанра
    PRIMARY KEY (genre_id)           -- Первичный ключ
) ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;