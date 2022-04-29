DROP TABLE IF EXISTS users;
CREATE TABLE users
(
    id         INT UNSIGNED primary key NOT NULL AUTO_INCREMENT,
    first_name VARCHAR(64)              NOT NULL,
    last_name  VARCHAR(64)              NOT NULL,
    age        integer unsigned         NOT NULL,
    email      VARCHAR(64)              NOT NULL,
    password   VARCHAR(255)             NOT NULL,
    created_at DATETIME                 NOT NULL,
    updated_at DATETIME                 NOT NULL,
    KEY index_users_on_email (email),
    KEY index_users_on_first_name_and_last_name (first_name, last_name)
) ENGINE = InnoDB
  DEFAULT character set = 'utf8mb4'
  collate = 'utf8mb4_general_ci' COMMENT = 'ユーザテーブル';