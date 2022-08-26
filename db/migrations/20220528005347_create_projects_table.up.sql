DROP TABLE IF EXISTS `projects`;
CREATE TABLE `projects`
(
    id                  char(36)                            NOT NULL comment 'ID',
    project_title       varchar(64)                         NOT NULL comment 'タイトル',
    project_description varchar(255)                        NOT NULL comment '説明',
    created_at          timestamp default current_timestamp NOT NULL comment '作成日時',
    updated_at          timestamp default current_timestamp NOT NULL comment '更新日時',
    PRIMARY KEY (id),
    KEY index_projects_on_project_title (project_title),
    KEY index_projects_on_created_at (created_at)
) ENGINE = InnoDB
  DEFAULT character
      set = 'utf8mb4'
  collate = 'utf8mb4_general_ci'
    COMMENT
        = 'プロジェクトテーブル';