-- ============================================================
-- 000001_init_schema
-- FortuneFinger 初始表结构：8 张表一次性建立。
--
-- 为什么是「一个文件」而不是「一张表一个文件」：
-- 这 8 张表是同一个意图——「项目诞生」。当前未定义数据库外键，
-- 也不存在「只回滚其中一张」的业务场景（没人会只删掉 dict 表
-- 而留着 user 表）。按变更意图拆分，而不是按表数量拆分。
--
-- 内容与 docs/schema.sql 保持一致，顺序也刻意保留，便于对照。
-- ============================================================

-- ---------- 系统用户表 ----------
CREATE TABLE `fortune_user` (
  `id` bigint unsigned NOT NULL COMMENT '主键ID',
  `username` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `nickname` varchar(64) NOT NULL DEFAULT '' COMMENT '昵称',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '密码哈希',
  `avatar` varchar(512) NOT NULL DEFAULT '' COMMENT '头像地址',
  `email` varchar(254) NOT NULL DEFAULT '' COMMENT '电子邮箱',
  `mobile` varchar(32) NOT NULL DEFAULT '' COMMENT '手机号码',
  `status` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '状态：1启用 2禁用',
  `created_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人',
  `updated_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新人',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除标识：非0已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fortune_user_username_deleted_at` (`username`, `deleted_at`),
  KEY `idx_fortune_user_deleted_at` (`deleted_at`),
  KEY `idx_fortune_user_email` (`email`),
  KEY `idx_fortune_user_mobile` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC COMMENT='系统用户表';

-- ---------- Casbin 权限策略表 ----------
CREATE TABLE `casbin_rule` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `ptype` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '策略类型：p为权限策略，g为角色关系',
    `v0` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '策略参数0',
    `v1` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '策略参数1',
    `v2` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '策略参数2',
    `v3` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '策略参数3',
    `v4` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '策略参数4',
    `v5` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '策略参数5',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_casbin_rule_policy` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`),
    KEY `idx_casbin_rule_v0` (`v0`),
    KEY `idx_casbin_rule_v1` (`v1`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC COMMENT = 'Casbin权限策略表';

-- ---------- 用户角色表 ----------
CREATE TABLE `fortune_role` (
  `id` bigint unsigned NOT NULL COMMENT '主键ID',
  `parent_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '父角色ID，0表示顶级',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '角色显示名',
  `code` varchar(100) NOT NULL COMMENT '角色编码',
  `status` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '状态：1启用 2禁用',
  `created_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人',
  `updated_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新人',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除标识：非0已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fortune_role_code_deleted_at` (`code`, `deleted_at`),
  KEY `idx_fortune_role_parent_deleted_status` (`parent_id`, `deleted_at`, `status`),
  KEY `idx_fortune_role_deleted_status` (`deleted_at`, `status`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC COMMENT='用户角色表';

-- ---------- 菜单表 ----------
CREATE TABLE `fortune_menu` (
  `id` bigint unsigned NOT NULL COMMENT '主键ID',
  `pid` bigint unsigned NOT NULL DEFAULT 0 COMMENT '父级ID',
  `types` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '节点类型：1目录 2菜单 3权限',
  `title` varchar(100) NOT NULL COMMENT '节点名称',
  `node` varchar(255) NOT NULL COMMENT '节点',
  `route` varchar(255) DEFAULT NULL COMMENT '路由',
  `code` varchar(100) DEFAULT NULL COMMENT '权限编码',
  `icon` varchar(100) NOT NULL DEFAULT '' COMMENT '图标',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `sort` int unsigned NOT NULL DEFAULT 0 COMMENT '排序',
  `keep_alive` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '页面缓存：1缓存 2不缓存',
  `status` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '状态：1启用 2禁用',
  `created_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人',
  `updated_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新人',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除标识',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fortune_menu_node_deleted_at` (`node`, `deleted_at`),
  UNIQUE KEY `uk_fortune_menu_route_deleted_at` (`route`, `deleted_at`),
  UNIQUE KEY `uk_fortune_menu_code_deleted_at` (`code`, `deleted_at`),
  KEY `idx_fortune_menu_pid_deleted_status_sort` (`pid`, `deleted_at`, `status`, `sort`),
  KEY `idx_fortune_menu_types_deleted_status` (`types`, `deleted_at`, `status`),
  KEY `idx_fortune_menu_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC COMMENT='菜单表';

-- ---------- 字典类型表 ----------
CREATE TABLE `fortune_dict_type` (
  `id` bigint unsigned NOT NULL COMMENT '主键ID',
  `name` varchar(100) NOT NULL COMMENT '字典名称',
  `code` varchar(100) NOT NULL COMMENT '字典编码',
  `is_system` tinyint unsigned NOT NULL DEFAULT 2 COMMENT '系统字典：1是 2否',
  `status` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '状态：1启用 2禁用',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人',
  `updated_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新人',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除标识',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fortune_dict_type_code_deleted_at` (`code`, `deleted_at`),
  KEY `idx_fortune_dict_type_deleted_status` (`deleted_at`, `status`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC COMMENT='字典类型表';

-- ---------- 字典数据表 ----------
CREATE TABLE `fortune_dict_data` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `dict_code` varchar(100) NOT NULL COMMENT '字典编码',
  `label` varchar(100) NOT NULL COMMENT '字典标签',
  `value` varchar(255) NOT NULL COMMENT '字典值',
  `color` varchar(32) NOT NULL DEFAULT 'default' COMMENT '展示颜色',
  `is_default` tinyint unsigned NOT NULL DEFAULT 2 COMMENT '默认选项：1是 2否',
  `sort` int unsigned NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '状态：1启用 2禁用',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人',
  `updated_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新人',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除标识',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fortune_dict_data_code_value_deleted_at` (`dict_code`, `value`, `deleted_at`),
  KEY `idx_fortune_dict_data_code_deleted_status_sort` (`dict_code`, `deleted_at`, `status`, `sort`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC COMMENT='字典数据表';

-- ---------- 组织结构表 ----------
CREATE TABLE `fortune_org` (
  `id` bigint unsigned NOT NULL COMMENT '主键ID',
  `parent_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '父级组织ID，0表示顶级组织',
  `name` varchar(100) NOT NULL COMMENT '组织名称',
  `code` varchar(100) NOT NULL COMMENT '组织编码',
  `types` tinyint unsigned NOT NULL COMMENT '组织类型：1集团 2公司 3部门',
  `sort` int unsigned NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '状态：1启用 2禁用',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人',
  `updated_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新人',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除标识：非0已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fortune_org_parent_code_deleted_at` (`parent_id`, `code`, `deleted_at`),
  UNIQUE KEY `uk_fortune_org_parent_name_deleted_at` (`parent_id`, `name`, `deleted_at`),
  KEY `idx_fortune_org_parent_deleted_status_sort` (`parent_id`, `deleted_at`, `status`, `sort`),
  KEY `idx_fortune_org_types_deleted_status` (`types`, `deleted_at`, `status`),
  KEY `idx_fortune_org_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC COMMENT='组织结构表';

-- ---------- 用户组织归属表 ----------
CREATE TABLE `fortune_org_user` (
  `id` bigint unsigned NOT NULL COMMENT '主键ID',
  `user_id` bigint unsigned NOT NULL COMMENT '系统用户ID',
  `company_id` bigint unsigned NOT NULL COMMENT '所属企业组织ID',
  `department_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '所属部门组织ID，0表示未分配部门',
  `relation_type` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '用户类型：1普通成员 2公司管理员 3集团管理员',
  `created_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人',
  `updated_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新人',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除标识：非0已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fortune_org_user_user_deleted_at` (`user_id`, `deleted_at`),
  KEY `idx_fortune_org_user_company_deleted_user` (`company_id`, `deleted_at`, `user_id`),
  KEY `idx_fortune_org_user_department_deleted_user` (`department_id`, `deleted_at`, `user_id`),
  KEY `idx_fortune_org_user_company_type_deleted` (`company_id`, `relation_type`, `deleted_at`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC COMMENT='用户组织归属表';

-- ---------- 对象存储记录表 ----------
CREATE TABLE `fortune_storage` (
  `id` bigint unsigned NOT NULL COMMENT '存储ID',
  `provider` varchar(20) NOT NULL COMMENT '存储提供方：minio、oss',
  `bucket` varchar(100) NOT NULL COMMENT '存储桶名称',
  `path` varchar(512) NOT NULL COMMENT '对象存储路径，不包含存储桶名称',
  `name` varchar(255) NOT NULL COMMENT '原始文件名',
  `ext` varchar(20) NOT NULL DEFAULT '' COMMENT '文件扩展名，不包含点号',
  `mime_type` varchar(100) NOT NULL COMMENT '文件MIME类型',
  `size` bigint unsigned NOT NULL DEFAULT 0 COMMENT '文件大小，单位字节',
  `hash` char(64) NOT NULL DEFAULT '' COMMENT '文件内容SHA-256摘要',
  `scene` varchar(32) NOT NULL COMMENT '使用场景：avatar、prize、image、video、event',
  `created_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人',
  `updated_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新人',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除标识：非0已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fortune_storage_provider_bucket_path_deleted_at` (`provider`, `bucket`, `path`, `deleted_at`),
  KEY `idx_fortune_storage_creator_deleted_created` (`created_by`, `deleted_at`, `created_at`),
  KEY `idx_fortune_storage_scene_deleted_created` (`scene`, `deleted_at`, `created_at`),
  KEY `idx_fortune_storage_deleted_created` (`deleted_at`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC COMMENT='对象存储记录表';
-- ---------- 奖品表 ----------
CREATE TABLE `fortune_prize` (
  `id` bigint unsigned NOT NULL COMMENT '奖品ID',
  `name` varchar(100) NOT NULL COMMENT '奖品名称',
  `storage_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '奖品图片存储ID',
  `org_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '所属组织ID',
  `level` varchar(100) NOT NULL DEFAULT 'participation' COMMENT '奖项等级',
  `types` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '奖品类型',
  `quantity` int unsigned NOT NULL DEFAULT 1 COMMENT '奖品数量',
  `status` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '状态：1启用 2禁用',
  `sort` int unsigned NOT NULL DEFAULT 0 COMMENT '排序',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '奖品说明',
  `created_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人',
  `updated_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新人',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除标识：非0已删除',
  PRIMARY KEY (`id`),
  KEY `idx_fortune_prize_level_deleted_status` (`level`, `deleted_at`, `status`),
  KEY `idx_fortune_prize_types_deleted_status` (`types`, `deleted_at`, `status`),
  KEY `idx_fortune_prize_deleted_status_sort` (`deleted_at`, `status`, `sort`, `id`),
  KEY `idx_fortune_prize_storage_deleted` (`storage_id`, `deleted_at`),
  KEY `idx_fortune_prize_org_deleted_status_sort` (`org_id`, `deleted_at`, `status`, `sort`, `id`),
  CONSTRAINT `chk_fortune_prize_quantity` CHECK (`quantity` >= 1)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC COMMENT='奖品表';

-- ---------- 系统设置表 ----------
CREATE TABLE `fortune_sys_config` (
  `id` bigint unsigned NOT NULL COMMENT '系统设置ID',
  `code` varchar(100) NOT NULL COMMENT '设置编码',
  `config` text NOT NULL COMMENT '设置内容',
  `created_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人',
  `updated_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新人',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fortune_sys_config_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC COMMENT='系统设置表';

-- ---------- 抽奖活动表 ----------
CREATE TABLE `fortune_activity` (
  `id` bigint unsigned NOT NULL COMMENT '活动ID',
  `code` char(8) NOT NULL COMMENT '活动编码',
  `name` varchar(100) NOT NULL COMMENT '活动名称',
  `org_id` bigint unsigned NOT NULL COMMENT '归属组织ID',
  `venue` varchar(255) NOT NULL DEFAULT '' COMMENT '活动场地',
  `activity_start_time` datetime(3) NOT NULL COMMENT '活动开始时间',
  `activity_end_time` datetime(3) NOT NULL COMMENT '活动结束时间',
  `checkin_start_time` datetime(3) NOT NULL COMMENT '签到开始时间',
  `checkin_end_time` datetime(3) NOT NULL COMMENT '签到结束时间',
  `checkin_template` varchar(50) NOT NULL DEFAULT 'invitation' COMMENT '签到模板',
  `status` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '状态：1启用 2禁用',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '活动说明',
  `created_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人',
  `updated_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新人',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除标识：非0已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fortune_activity_code` (`code`),
  KEY `idx_fortune_activity_org_deleted_status_start` (`org_id`, `deleted_at`, `status`, `activity_start_time`, `id`),
  KEY `idx_fortune_activity_deleted_status_start` (`deleted_at`, `status`, `activity_start_time`, `id`),
  CONSTRAINT `chk_fortune_activity_time` CHECK (`activity_end_time` > `activity_start_time`),
  CONSTRAINT `chk_fortune_activity_checkin_time` CHECK (`checkin_end_time` > `checkin_start_time`),
  CONSTRAINT `chk_fortune_activity_checkin_end` CHECK (`checkin_end_time` <= `activity_end_time`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC COMMENT='抽奖活动表';

-- ---------- 抽奖轮次表 ----------
CREATE TABLE `fortune_round` (
  `id` bigint unsigned NOT NULL COMMENT '轮次ID',
  `activity_id` bigint unsigned NOT NULL COMMENT '活动ID',
  `prize_id` bigint unsigned NOT NULL COMMENT '奖品ID',
  `name` varchar(100) NOT NULL COMMENT '轮次名称',
  `draw_date` date NOT NULL COMMENT '开奖日期',
  `draw_time` time NOT NULL COMMENT '开奖时间',
  `winner_count` int unsigned NOT NULL DEFAULT 1 COMMENT '中奖人数',
  `created_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人',
  `updated_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新人',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除标识：非0已删除',
  PRIMARY KEY (`id`),
  KEY `idx_fortune_round_activity_deleted_draw` (`activity_id`, `deleted_at`, `draw_date`, `draw_time`, `id`),
  KEY `idx_fortune_round_prize_deleted_draw` (`prize_id`, `deleted_at`, `draw_date`, `draw_time`, `id`),
  CONSTRAINT `chk_fortune_round_winner_count` CHECK (`winner_count` >= 1)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC COMMENT='抽奖轮次表';

-- ---------- 活动签到表 ----------
CREATE TABLE `fortune_checkin` (
  `id` bigint unsigned NOT NULL COMMENT '签到ID',
  `activity_id` bigint unsigned NOT NULL COMMENT '活动ID',
  `org_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '所属企业ID',
  `roster_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '签到名单ID，0表示自由签到',
  `name` varchar(100) NOT NULL COMMENT '参与者姓名',
  `mobile` varchar(32) NOT NULL COMMENT '手机号码',
  `email` varchar(254) NOT NULL DEFAULT '' COMMENT '电子邮箱',
  `identity_hash` char(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL COMMENT '身份摘要',
  `avatar` bigint unsigned NOT NULL DEFAULT 0 COMMENT '头像存储ID',
  `checked_in_at` datetime(3) NOT NULL COMMENT '签到时间',
  `created_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人',
  `updated_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新人',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除标识：非0已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_fortune_checkin_activity_identity_deleted` (`activity_id`, `identity_hash`, `deleted_at`),
  KEY `idx_fortune_checkin_activity_deleted_checked` (`activity_id`, `deleted_at`, `checked_in_at`, `id`),
  KEY `idx_fortune_checkin_org_activity_deleted` (`org_id`, `activity_id`, `deleted_at`),
  KEY `idx_fortune_checkin_roster_deleted` (`roster_id`, `deleted_at`)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC COMMENT='活动签到表';

-- ---------- 签到指纹表 ----------
CREATE TABLE `fortune_checkin_fingerprint` (
  `id` bigint unsigned NOT NULL COMMENT '指纹记录ID',
  `activity_id` bigint unsigned NOT NULL COMMENT '活动ID',
  `checkin_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '签到ID，0表示尚未完成签到',
  `identity_hash` char(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '' COMMENT '身份摘要',
  `session_hash` char(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '' COMMENT '参与者会话摘要',
  `browser_instance_hash` char(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '' COMMENT '浏览器实例摘要',
  `fingerprint_hash` char(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '' COMMENT '综合设备指纹摘要',
  `fingerprint_version` smallint unsigned NOT NULL DEFAULT 1 COMMENT '指纹算法版本',
  `fingerprintjs_hash` char(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '' COMMENT 'FingerprintJS摘要',
  `canvas_hash` char(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '' COMMENT 'Canvas特征摘要',
  `webgl_hash` char(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '' COMMENT 'WebGL特征摘要',
  `audio_hash` char(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '' COMMENT 'AudioContext特征摘要',
  `screen_hash` char(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '' COMMENT '屏幕特征摘要',
  `user_agent_hash` char(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '' COMMENT '用户代理摘要',
  `client_hints_hash` char(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '' COMMENT '客户端提示摘要',
  `languages_hash` char(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '' COMMENT '语言列表摘要',
  `screen_width` int unsigned NOT NULL DEFAULT 0 COMMENT '屏幕宽度',
  `screen_height` int unsigned NOT NULL DEFAULT 0 COMMENT '屏幕高度',
  `avail_width` int unsigned NOT NULL DEFAULT 0 COMMENT '屏幕可用宽度',
  `avail_height` int unsigned NOT NULL DEFAULT 0 COMMENT '屏幕可用高度',
  `device_pixel_ratio` decimal(6,3) unsigned NOT NULL DEFAULT 0 COMMENT '设备像素比',
  `color_depth` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '颜色深度',
  `timezone` varchar(64) NOT NULL DEFAULT '' COMMENT '浏览器时区',
  `platform` varchar(64) NOT NULL DEFAULT '' COMMENT '设备平台',
  `browser` varchar(64) NOT NULL DEFAULT '' COMMENT '浏览器类型',
  `hardware_concurrency` smallint unsigned NOT NULL DEFAULT 0 COMMENT '逻辑处理器数量',
  `device_memory` decimal(6,2) unsigned NOT NULL DEFAULT 0 COMMENT '设备内存档位',
  `max_touch_points` smallint unsigned NOT NULL DEFAULT 0 COMMENT '最大触控点数量',
  `touch_supported` tinyint unsigned NOT NULL DEFAULT 2 COMMENT '触控支持：1是 2否',
  `ip_hash` char(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '' COMMENT '请求IP摘要',
  `components` json NOT NULL COMMENT '其他指纹分项摘要',
  `confidence` decimal(5,4) unsigned NOT NULL DEFAULT 0 COMMENT '指纹置信度',
  `similarity_score` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '设备相似度',
  `risk_reason` varchar(500) NOT NULL DEFAULT '' COMMENT '风险原因',
  `decision` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '处理结果：0待判断 1通过 2待验证 3拒绝 4验证通过',
  `created_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人',
  `updated_by` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新人',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除标识：非0已删除',
  PRIMARY KEY (`id`),
  KEY `idx_fortune_checkin_fp_activity_browser_deleted` (`activity_id`, `browser_instance_hash`, `deleted_at`),
  KEY `idx_fortune_checkin_fp_activity_fingerprint_deleted` (`activity_id`, `fingerprint_hash`, `deleted_at`),
  KEY `idx_fortune_checkin_fp_activity_visitor_deleted` (`activity_id`, `fingerprintjs_hash`, `deleted_at`),
  KEY `idx_fortune_checkin_fp_activity_identity_deleted` (`activity_id`, `identity_hash`, `deleted_at`),
  KEY `idx_fortune_checkin_fp_activity_session_deleted` (`activity_id`, `session_hash`, `deleted_at`),
  KEY `idx_fortune_checkin_fp_checkin_deleted` (`checkin_id`, `deleted_at`),
  KEY `idx_fortune_checkin_fp_activity_decision_created` (`activity_id`, `decision`, `deleted_at`, `created_at`),
  CONSTRAINT `chk_fortune_checkin_fingerprint_confidence` CHECK (`confidence` <= 1),
  CONSTRAINT `chk_fortune_checkin_fingerprint_similarity` CHECK (`similarity_score` <= 100),
  CONSTRAINT `chk_fortune_checkin_fingerprint_decision` CHECK (`decision` <= 4)
) ENGINE=InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC COMMENT='签到指纹表';
