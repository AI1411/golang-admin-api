DROP TABLE IF EXISTS `subscription_members`;
CREATE TABLE `subscription_members`
(
    id                char(36)                            NOT NULL comment 'ID',
    user_id           char(36)                            NOT NULL comment 'ユーザID',
    member_status     varchar(64)                         NOT NULL comment '会員ステータス',
    member_start_date timestamp default current_timestamp NOT NULL comment '会員開始日',
    member_end_date   timestamp default current_timestamp NOT NULL comment '会員開始日',
    created_at        timestamp default current_timestamp NOT NULL comment '作成日時',
    updated_at        timestamp default current_timestamp NOT NULL comment '更新日時'
) ENGINE = InnoDB
  DEFAULT character
      set = 'utf8mb4'
  collate = 'utf8mb4_general_ci'
    COMMENT
        = '定額課金ユーザ';