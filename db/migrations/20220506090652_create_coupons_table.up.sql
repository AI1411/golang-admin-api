DROP TABLE IF EXISTS coupons;
CREATE TABLE coupons
(
    id                  char(36) primary key                NOT NULL comment 'ID',
    title               varchar(64)                         NOT NULL comment 'タイトル',
    remarks             text                                NULL comment '説明',
    discount_amount     integer                             NULL comment '値引額',
    discount_rate       integer                             NULL comment '割引率',
    max_discount_amount integer                             NULL comment '最大値引額',
    use_start_at        timestamp                           NOT NULL comment '利用開始日',
    use_end_at          timestamp                           NOT NULL comment '利用終了日',
    public_start_at     timestamp                           NOT NULL comment '公開開始日',
    public_end_at       timestamp                           NOT NULL comment '公開終了日',
    is_public           boolean   default false             NOT NULL comment '公開フラグ',
    is_premium          boolean   default false             NOT NULL comment 'プレミアムフラグ',
    created_at          timestamp default current_timestamp NOT NULL comment '作成日時',
    updated_at          timestamp default current_timestamp NOT NULL comment '更新日時',
    KEY index_products_on_id (id),
    KEY index_products_on_title (title),
    KEY index_products_on_is_public (is_public),
    KEY index_products_on_discount_amount (discount_amount),
    KEY index_products_on_discount_rate (discount_rate),
    KEY index_products_on_max_discount_amount (max_discount_amount)
) ENGINE = InnoDB
  DEFAULT character
      set = 'utf8mb4'
  collate = 'utf8mb4_general_ci'
    COMMENT
        = 'クーポンテーブル';