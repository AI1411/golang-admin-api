DROP TABLE IF EXISTS `group_user`;
CREATE TABLE `group_user`
(
    id         integer auto_increment primary key  NOT NULL comment 'ID',
    group_id   char(36)                            NOT NULL comment 'グループID',
    user_id    char(36)                            NOT NULL comment 'ユーザID',
    created_at timestamp default current_timestamp NOT NULL comment '作成日時',
    updated_at timestamp default current_timestamp NOT NULL comment '更新日時',
    KEY index_products_on_group_id (group_id),
    KEY index_products_on_user_id (user_id)
) ENGINE = InnoDB
  DEFAULT character
      set = 'utf8mb4'
  collate = 'utf8mb4_general_ci'
    COMMENT
        = 'グループとユーザの中間テーブル';