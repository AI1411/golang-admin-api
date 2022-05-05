DROP TABLE IF EXISTS todos;
CREATE TABLE todos
(
    id         INT UNSIGNED primary key           NOT NULL AUTO_INCREMENT,
    title      VARCHAR(64)                        NOT NULL comment 'タイトル',
    body       TEXT                               NOT NULL comment '本文',
    status     varchar(64)                        NOT NULL comment 'ステータス',
    user_id    char(36)                           NOT NULL comment 'ユーザーID',
    created_at DATETIME default current_timestamp NOT NULL comment '作成日時',
    updated_at DATETIME default current_timestamp NOT NULL comment '更新日時',
    KEY user_id_idx (user_id),
    KEY title_idx (title),
    KEY status_idx (status)
) ENGINE = InnoDB
  DEFAULT character set = 'utf8mb4'
  collate = 'utf8mb4_general_ci' COMMENT = 'TODOテーブル';