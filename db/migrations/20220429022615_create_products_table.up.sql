DROP TABLE IF EXISTS products;
CREATE TABLE products
(
    id           char(36) primary key                   NOT NULL comment '商品ID',
    product_name varchar(64)                            NOT NULL comment '商品名',
    price        int unsigned                           NOT NULL comment '商品価格',
    remarks      varchar(255) default NULL comment '商品備考',
    quantity     int unsigned default 1                 NOT NULL comment '商品数量',
    created_at   timestamp    default current_timestamp NOT NULL comment '作成日時',
    updated_at   timestamp    default current_timestamp NOT NULL comment '更新日時',
    KEY index_products_on_id (id),
    KEY index_products_on_name (product_name),
    KEY index_products_on_price (price)
) ENGINE = InnoDB
  DEFAULT character set = 'utf8mb4'
  collate = 'utf8mb4_general_ci' COMMENT = '商品テーブル';