DROP TABLE IF EXISTS coupon_user;
CREATE TABLE coupon_user
(
    id         int auto_increment primary key               NOT NULL comment 'ID',
    coupon_id  char(36)                                     NOT NULL comment 'クーポンID',
    user_id    char(36)                                     NOT NULL comment 'ユーザID',
    use_count  mediumint unsigned default 0                 NOT NULL comment '使用回数',
    created_at timestamp          default current_timestamp NOT NULL comment '作成日時',
    updated_at timestamp          default current_timestamp NOT NULL comment '更新日時',
    KEY index_products_on_id (id),
    KEY index_products_on_coupon_id (coupon_id),
    KEY index_products_on_user_id (user_id)
) ENGINE = InnoDB
  DEFAULT character
      set = 'utf8mb4'
  collate = 'utf8mb4_general_ci'
    COMMENT
        = 'クーポン紐付けテーブル';