CREATE TABLE `tb_user`
(
    `id`         bigint       NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username`  varchar(32)  NOT NULL COMMENT '用户名',
    `avatar`     varchar(255)          DEFAULT NULL COMMENT '头像',
    `sex`        tinyint(2)            DEFAULT '1' COMMENT '性别(0:女,1:男)',
    `password`   varchar(128) NOT NULL COMMENT '密码(加密存储)',
    `email`      varchar(100)          DEFAULT NULL COMMENT '邮箱',
    `phone`      varchar(20)           DEFAULT NULL COMMENT '手机号',
    `status`     tinyint(2)   NOT NULL DEFAULT '1' COMMENT '状态(0:禁用,1:启用)',
    `created_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_email` (`email`),
    UNIQUE KEY `idx_phone` (`phone`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='用户表';

CREATE TABLE `tb_role` (
                           `id` bigint NOT NULL AUTO_INCREMENT COMMENT '角色ID',
                           `role_name` varchar(32) NOT NULL COMMENT '角色名',
                           `status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '状态(0:禁用,1:启用)',
                           `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                           PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';
-- 插入基础角色
INSERT INTO `tb_role` (`role_name`, `status`) VALUES
('普通用户', 1),
('商户', 1);
                                                                           

CREATE TABLE `tb_user_role` (
                                `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
                                `user_id` bigint NOT NULL COMMENT '用户ID',
                                `role_id` bigint NOT NULL COMMENT '角色ID',
                                `status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '状态(0:禁用,1:启用)',
                                `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                PRIMARY KEY (`id`),
                                UNIQUE KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

CREATE TABLE `tb_product` (
                              `id` bigint NOT NULL AUTO_INCREMENT COMMENT '商品ID',
                              `name` varchar(64) NOT NULL COMMENT '商品名称',
                              `description` text COMMENT '商品描述',
                              `picture` varchar(255) COMMENT '商品图片',
                              `price` decimal(10,2) NOT NULL COMMENT '商品价格',
                              `stock` int NOT NULL DEFAULT '0' COMMENT '库存数量',
                              `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态(0:下架,1:上架)',
                              `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                              `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                              PRIMARY KEY (`id`),
                              KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';

CREATE TABLE `tb_shopping_cart` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '购物车ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='购物车表';

CREATE TABLE `tb_cart_item` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '购物车项ID',
  `cart_id` bigint NOT NULL COMMENT '购物车ID',
  `product_id` bigint NOT NULL COMMENT '商品ID',
  `quantity` int NOT NULL DEFAULT '1' COMMENT '商品数量',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_cart_id` (`cart_id`),
  KEY `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='购物车项表';