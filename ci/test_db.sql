-- date: 2026-01-06 pm
-- author: wendisx
-- with: golang
-- desc: test some basic sql functions

-- generate by ai, exactly i just don't want to do it by myself, this is the easy work cost time a lot.
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- user_basic
DROP TABLE IF EXISTS `user_basic`;
CREATE TABLE `user_basic` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `extern_id` binary(16) NOT NULL COMMENT '拓展ID（UUID二进制存储）',
  `user_name` varchar(20) NOT NULL COMMENT '用户名',
  `user_password` varchar(64) NOT NULL COMMENT '用户密码（建议存储加密后的值，长度适配MD5/SHA256）',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP() COMMENT '更新时间',
  `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除：0-未删除，1-已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_name` (`user_name`),
  UNIQUE KEY `uk_extern_id` (`extern_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户基础信息表';

-- user_detail
DROP TABLE IF EXISTS `user_detail`;
CREATE TABLE `user_detail` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` bigint unsigned NOT NULL COMMENT '关联user_basic的主键ID',
  `nickname` varchar(30) DEFAULT '' COMMENT '用户昵称',
  `phone` varchar(11) DEFAULT '' COMMENT '手机号',
  `email` varchar(50) DEFAULT '' COMMENT '邮箱',
  `avatar` varchar(255) DEFAULT '' COMMENT '头像地址',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP() COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id` (`user_id`),
  CONSTRAINT `fk_user_detail_user_id` FOREIGN KEY (`user_id`) REFERENCES `user_basic` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户详情信息表';

SET FOREIGN_KEY_CHECKS = 1;

-- some sql test can be here.