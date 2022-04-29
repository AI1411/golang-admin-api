DROP TABLE IF EXISTS todos;
CREATE TABLE todos
(
    id         INT UNSIGNED primary key NOT NULL AUTO_INCREMENT,
    title      VARCHAR(64)              NOT NULL,
    body       TEXT                     NOT NULL,
    status     varchar(64)              NOT NULL,
    user_id    INT UNSIGNED             NOT NULL,
    created_at DATETIME                 NOT NULL,
    updated_at DATETIME                 NOT NULL,
    KEY user_id_idx (user_id),
    KEY title_idx (title),
    KEY status_idx (status)
) ENGINE = InnoDB
  DEFAULT character set = 'utf8mb4'
  collate = 'utf8mb4_general_ci' COMMENT = 'TODOテーブル';