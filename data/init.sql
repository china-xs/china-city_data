CREATE TABLE `m_district` (
  `id` int NOT NULL AUTO_INCREMENT,
  `p_id` int unsigned NOT NULL DEFAULT '0' COMMENT '上级ID',
  `p_ids` varchar(64) NOT NULL DEFAULT '' COMMENT '上级IDs',
  `name` varchar(50) DEFAULT NULL COMMENT '名称',
  `level` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '层级',
  `merger` varchar(128) NOT NULL DEFAULT '' COMMENT '地区全名',
  `code` varchar(48) NOT NULL DEFAULT '' COMMENT '统计用区划代码',
  `order_num` int unsigned NOT NULL DEFAULT '0' COMMENT '排序 order_num asc,created_at asc ',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=678617 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='2020-国家统计局-城乡数据';