CREATE TABLE user
(
    user_id   BINARY(16) NOT NULL,          -- UUID пользователя
    avatar_id BINARY(16),                   -- UUID аватара (может быть NULL)
    login     VARCHAR(255) UNIQUE NOT NULL, -- Уникальный логин
    role      SMALLINT,                     -- Роль пользователя
    password  VARCHAR(255)        NOT NULL, -- Пароль
    about_me  TEXT,                         -- Описание о себе
    PRIMARY KEY (user_id)                   -- Первичный ключ
) ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;