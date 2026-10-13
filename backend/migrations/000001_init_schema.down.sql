-- ============================================================
-- 000001_init_schema 的回滚
-- 精确撤销 up：删除 up 中建立的全部 8 张表。
-- 顺序与 up 相反（先建的后删）；当前无外键，顺序不影响执行，
-- 但保持逆序是规范做法，将来加外键时不会出错。
--
-- 注意：DROP TABLE 会连同数据一起删除，且不可恢复。
-- ============================================================

DROP TABLE IF EXISTS `fortune_org_user`;
DROP TABLE IF EXISTS `fortune_org`;
DROP TABLE IF EXISTS `fortune_dict_data`;
DROP TABLE IF EXISTS `fortune_dict_type`;
DROP TABLE IF EXISTS `fortune_menu`;
DROP TABLE IF EXISTS `fortune_role`;
DROP TABLE IF EXISTS `casbin_rule`;
DROP TABLE IF EXISTS `fortune_user`;
