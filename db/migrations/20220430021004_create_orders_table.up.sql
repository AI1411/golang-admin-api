DROP TABLE IF EXISTS orders;
CREATE TABLE orders
(
    id           char(36) primary key                NOT NULL comment '注文ID',
    user_id      char(36)                            NOT NULL comment 'ユーザID',
    quantity     varchar(64)                         NOT NULL comment '数量',
    total_price  mediumint unsigned                  NOT NULL comment '合計金額',
    order_status varchar(64)                         NULL comment '注文ステータス',
    remarks      text                                NOT NULL comment '注文備考',
    created_at   timestamp default current_timestamp NOT NULL comment '作成日時',
    updated_at   timestamp default current_timestamp NOT NULL comment '更新日時',
    KEY index_products_on_id (id),
    KEY index_products_on_user_id (user_id),
    KEY index_products_on_order_status (order_status)
) ENGINE = InnoDB
  DEFAULT character set = 'utf8mb4'
  collate = 'utf8mb4_general_ci' COMMENT = '注文テーブル';