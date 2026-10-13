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
