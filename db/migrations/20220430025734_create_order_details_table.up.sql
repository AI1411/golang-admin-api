DROP TABLE IF EXISTS order_details;
CREATE TABLE order_details
(
    id         char(36) primary key                NOT NULL comment '注文詳細ID',
    order_id   char(36)                            NOT NULL comment '注文ID',
    product_id char(36)                            NOT NULL comment '商品ID',
    quantity   varchar(64)                         NOT NULL comment '数量',
    price      integer unsigned                    NOT NULL comment '価格',
    created_at timestamp default current_timestamp NOT NULL comment '作成日時',
    updated_at timestamp default current_timestamp NOT NULL comment '更新日時',
    KEY index_products_on_id (id),
    KEY index_products_on_order_id (order_id),
    KEY index_products_on_product_id (product_id)
) ENGINE = InnoDB
  DEFAULT character set = 'utf8mb4'
  collate = 'utf8mb4_general_ci' COMMENT = '注文詳細テーブル';