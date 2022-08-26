DROP TABLE IF EXISTS `epics`;
CREATE TABLE `epics`
(
    id               integer auto_increment primary key  NOT NULL comment 'ID',
    is_open          boolean   default false             NOT NULL DEFAULT true comment '解放フラグ',
    author_id        char(36)                            NULL comment '作成者ID',
    epic_title       varchar(64)                         NULL comment 'タイトル',
    epic_description varchar(255)                        NULL comment '説明',
    label            varchar(64)                         NULL comment 'ラベル',
    milestone_id     char(36)                            NULL comment 'マイルストーンID',
    assignee_id      char(36)                            NULL comment '担当者ID',
    project_id       char(36)                            NULL comment 'プロジェクトID',
    created_at       timestamp default current_timestamp NOT NULL comment '作成日時',
    updated_at       timestamp default current_timestamp NOT NULL comment '更新日時',
    KEY index_products_on_milestone_id (milestone_id),
    KEY index_products_on_author_id (author_id),
    KEY index_products_on_assignee_id (assignee_id)
) ENGINE = InnoDB
  DEFAULT character
      set = 'utf8mb4'
  collate = 'utf8mb4_general_ci'
    COMMENT
        = 'エピックテーブル';