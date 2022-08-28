DROP TABLE IF EXISTS `issues`;
CREATE TABLE `issues`
(
    id           char(36) primary key                  NOT NULL comment 'ID',
    title        varchar(64)                           NOT NULL comment 'Issue title',
    description  varchar(255)                          NOT NULL comment 'Issue description',
    user_id      char(36)                              NOT NULL comment 'ユーザID',
    milestone_id char(36)                              NOT NULL comment 'マイルストーンID',
    issue_status varchar(64) default 'waiting'         NOT NULL comment 'Issueステータス',
    created_at   timestamp   default current_timestamp NOT NULL comment '作成日時',
    updated_at   timestamp   default current_timestamp NOT NULL comment '更新日時',
    KEY index_user_id (user_id),
    KEY index_issue_title (title),
    KEY index_issue_status (issue_status)
) ENGINE = InnoDB
  DEFAULT character
      set = 'utf8mb4'
  collate = 'utf8mb4_general_ci'
    COMMENT
        = 'Issueテーブル';