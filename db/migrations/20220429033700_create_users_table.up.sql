DROP TABLE IF EXISTS users;
CREATE TABLE users
(
    id         int primary key auto_increment     NOT NULL comment 'ID',
    first_name VARCHAR(64)                        NOT NULL comment '名前',
    last_name  VARCHAR(64)                        NOT NULL comment '姓',
    age        integer unsigned                   NOT NULL comment '年齢',
    email      VARCHAR(64)                        NOT NULL comment 'メールアドレス',
    password   VARCHAR(255)                       NOT NULL comment 'パスワード',
    created_at DATETIME default current_timestamp NOT NULL comment '作成日時',
    updated_at DATETIME default current_timestamp NOT NULL comment '更新日時',
    KEY index_users_on_email (email),
    KEY index_users_on_first_name_and_last_name (first_name, last_name)
) ENGINE = InnoDB
  DEFAULT character set = 'utf8mb4'
  collate = 'utf8mb4_general_ci' COMMENT = 'ユーザテーブル';