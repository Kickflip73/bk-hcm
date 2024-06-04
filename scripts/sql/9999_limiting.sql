-- hcm.limit_rule definition

CREATE TABLE `limit_rule` (
                              `id` bigint NOT NULL AUTO_INCREMENT COMMENT '限流规则标识id',
                              `rule_name` varchar(255) DEFAULT NULL COMMENT '限流规则名称',
                              `scene` varchar(255) DEFAULT NULL COMMENT '使用场景',
                              `account` varchar(255) DEFAULT NULL COMMENT '账号',
                              `identify` varchar(255) DEFAULT NULL COMMENT '标识',
                              `max_limit` bigint DEFAULT NULL COMMENT '最大限制',
                              `windows_size` int DEFAULT NULL COMMENT '窗口大小',
                              `reject_policy` varchar(255) DEFAULT NULL COMMENT '拒绝策略',
                              `retry_interval` int DEFAULT NULL COMMENT '重试间隔',
                              `retry_max_timeout` int DEFAULT NULL COMMENT '最大超时',
                              `retry_max_count` int DEFAULT NULL COMMENT '最大重试次数',
                              `deny_all` tinyint(1) DEFAULT NULL COMMENT '是否拒绝全部',
                              `enabled` tinyint(1) DEFAULT NULL COMMENT '是否启用',
                              `creator` varchar(255) DEFAULT NULL COMMENT '创建者',
                              `reviser` varchar(255) DEFAULT NULL COMMENT '更新者',
                              `created_at` datetime DEFAULT NULL COMMENT '创建时间',
                              `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
                              PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;