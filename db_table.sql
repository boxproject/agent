SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for T_HASH
-- ----------------------------
DROP TABLE IF EXISTS `T_HASH`;
CREATE TABLE `T_HASH` (
  `hash` varchar(100) NOT NULL COMMENT '模版id',
  `app_id` varchar(20) NOT NULL DEFAULT '' COMMENT '员工id',
  `captain_id` varchar(60) DEFAULT NULL COMMENT '私钥id',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '模板名称',
  `flow` text NOT NULL COMMENT '原始数据',
  `sign` text NOT NULL COMMENT '签名',
  `status` varchar(2) NOT NULL DEFAULT '' COMMENT '状态',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for T_HASH_OPERATE
-- ----------------------------
DROP TABLE IF EXISTS `T_HASH_OPERATE`;
CREATE TABLE `T_HASH_OPERATE` (
  `id` int(20) NOT NULL AUTO_INCREMENT COMMENT 'hash序列',
  `app_id` varchar(20) NOT NULL DEFAULT '' COMMENT '私钥id',
  `type` varchar(20) NOT NULL DEFAULT '' COMMENT '操作类型',
  `hash` varchar(100) NOT NULL DEFAULT '' COMMENT 'hash',
  `option` varchar(255) NOT NULL DEFAULT '' COMMENT '审批状态',
  `sign` text NOT NULL COMMENT '签名',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `opinion` varchar(255) NOT NULL DEFAULT '' COMMENT '操作意见',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=132 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for T_REGIST
-- ----------------------------
DROP TABLE IF EXISTS `T_REGIST`;
CREATE TABLE `T_REGIST` (
  `reg_id` varchar(100) NOT NULL COMMENT '注册id',
  `applyer_id` varchar(100) NOT NULL DEFAULT '' COMMENT '员工id',
  `captain_id` varchar(100) NOT NULL DEFAULT '' COMMENT '私钥id',
  `applyer_account` varchar(100) NOT NULL DEFAULT '' COMMENT '申请人账号',
  `msg` text NOT NULL COMMENT '信息',
  `consent` varchar(255) NOT NULL DEFAULT '',
  `cipher_text` varchar(255) NOT NULL DEFAULT '',
  `status` varchar(255) NOT NULL DEFAULT '' COMMENT '状态',
  `pub_key` text NOT NULL COMMENT '公钥',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`reg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for T_REQ_LOG
-- ----------------------------
DROP TABLE IF EXISTS `T_REQ_LOG`;
CREATE TABLE `T_REQ_LOG` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `req_type` varchar(255) NOT NULL DEFAULT '',
  `transfer_type` varchar(255) NOT NULL DEFAULT '',
  `block_number` bigint(20) NOT NULL DEFAULT '0',
  `hash` varchar(255) NOT NULL DEFAULT '',
  `wd_hash` varchar(255) NOT NULL DEFAULT '',
  `tx_hash` varchar(255) NOT NULL DEFAULT '',
  `amount` varchar(255) NOT NULL DEFAULT '',
  `fee` varchar(255) NOT NULL DEFAULT '',
  `from` varchar(255) NOT NULL DEFAULT '',
  `to` varchar(255) NOT NULL DEFAULT '',
  `category` bigint(20) NOT NULL DEFAULT '0',
  `content` varchar(255) NOT NULL DEFAULT '',
  `status` varchar(255) NOT NULL DEFAULT '',
  `apply_time` datetime NOT NULL,
  `create_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for T_WITHDRAW
-- ----------------------------
DROP TABLE IF EXISTS `T_WITHDRAW`;
CREATE TABLE `T_WITHDRAW` (
  `wd_hash` varchar(100) NOT NULL COMMENT '提现id',
  `hash` varchar(100) NOT NULL DEFAULT '' COMMENT '审批id',
  `to` varchar(60) NOT NULL DEFAULT '' COMMENT '转账地址',
  `amount` varchar(40) NOT NULL DEFAULT '' COMMENT '转账金额',
  `fee` varchar(40) NOT NULL DEFAULT '' COMMENT '手续费',
  `category` int(3) NOT NULL DEFAULT '0' COMMENT '币种编号',
  `flow` text NOT NULL COMMENT '原始数据',
  `sign` text NOT NULL COMMENT '签名',
  `status` varchar(2) NOT NULL DEFAULT '' COMMENT '提现状态',
  `app_id` varchar(20) NOT NULL DEFAULT '' COMMENT '员工id',
  `wd_flow` text NOT NULL,
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`wd_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;
