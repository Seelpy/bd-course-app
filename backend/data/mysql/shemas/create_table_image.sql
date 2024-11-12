CREATE TABLE images
(
    image_id BINARY(16) NOT NULL,   -- UUID изображения
    path     VARCHAR(255) NOT NULL, -- Путь к изображению
    PRIMARY KEY (image_id)          -- Первичный ключ
) ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;