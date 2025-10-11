CREATE TABLE `users` (
                         `id` VARCHAR(24) NOT NULL DEFAULT '' COMMENT '用户ID',
                         `avatar` VARCHAR(255) NOT NULL DEFAULT 'https://gw.alipayobjects.com/zos/rmsportal/BiazfanxmamNRoxxVxka.png' COMMENT '头像',
                         `nickname` VARCHAR(24) NOT NULL DEFAULT '' COMMENT '昵称',
                         `phone` VARCHAR(24) NOT NULL DEFAULT '' COMMENT '手机号',
                         `email` VARCHAR(24) DEFAULT NULL COMMENT '邮箱',
                         `password` VARCHAR(191) DEFAULT NULL COMMENT '密码',
                         `status` TINYINT(1) DEFAULT 0 COMMENT '状态 0-正常 1-禁用',
                         `sex` TINYINT(1) DEFAULT 0 COMMENT '性别 0-未知 1-男 2-女',
                         `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                         `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                         PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
