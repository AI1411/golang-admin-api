DROP TABLE IF EXISTS `milestones`;
CREATE TABLE `milestones`
(
    id                    char(36)                            NOT NULL comment 'ID',
    milestone_title       varchar(64)                         NOT NULL comment 'タイトル',
    milestone_description varchar(255)                        NOT NULL comment '説明',
    project_id            char(36)                            NOT NULL comment 'プロジェクトID',
    created_at            timestamp default current_timestamp NOT NULL comment '作成日時',
    updated_at            timestamp default current_timestamp NOT NULL comment '更新日時'
) ENGINE = InnoDB
  DEFAULT character
      set = 'utf8mb4'
  collate = 'utf8mb4_general_ci'
    COMMENT
        = 'マイルストーンテーブル';