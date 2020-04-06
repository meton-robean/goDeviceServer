/*
Navicat MySQL Data Transfer

Source Server         : 192.168.50.128
Source Server Version : 50560
Source Host           : 192.168.50.129:3306
Source Database       : outside_help

Target Server Type    : MYSQL
Target Server Version : 50560
File Encoding         : 65001

Date: 2018-06-20 22:08:02
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for t_agent_history_info
-- ----------------------------
DROP TABLE IF EXISTS `t_agent_history_info`;
CREATE TABLE `t_agent_history_info` (
  `id` int(4) NOT NULL AUTO_INCREMENT COMMENT '操作用户id',
  `user_id` int(4) NOT NULL DEFAULT '0' COMMENT '操作人的id',
  `msg` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '操作信息',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '入库时间戳',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='操作日志表';

-- ----------------------------
-- Records of t_agent_history_info
-- ----------------------------

-- ----------------------------
-- Table structure for t_agent_order_info
-- ----------------------------
DROP TABLE IF EXISTS `t_agent_order_info`;
CREATE TABLE `t_agent_order_info` (
  `id` int(4) NOT NULL AUTO_INCREMENT,
  `order_id` varchar(32) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '订单号',
  `agent_id` int(4) NOT NULL DEFAULT '0' COMMENT '关联代理商的id',
  `product_id` int(4) NOT NULL DEFAULT '0' COMMENT '产品id',
  `num` int(4) NOT NULL DEFAULT '0' COMMENT '订单数量',
  `status` int(2) NOT NULL DEFAULT '0' COMMENT '订单是否被确认 0 未确认 1 确认',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '订单时间戳',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='代理商订单操作';

-- ----------------------------
-- Records of t_agent_order_info
-- ----------------------------
INSERT INTO `t_agent_order_info` VALUES ('4', 'YZ2018053020540168', '25', '2', '5', '0', '1527684849');
INSERT INTO `t_agent_order_info` VALUES ('2', 'YZ2018051422040423', '1', '2', '2', '1', '1527609740');
INSERT INTO `t_agent_order_info` VALUES ('3', 'YZ2018053018574010', '19', '2', '20', '0', '1527677880');

-- ----------------------------
-- Table structure for t_device_bind_info
-- ----------------------------
DROP TABLE IF EXISTS `t_device_bind_info`;
CREATE TABLE `t_device_bind_info` (
  `id` int(4) NOT NULL AUTO_INCREMENT,
  `device_id` varchar(32) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '设备id，唯一的 md5',
  `roomnu` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '房间号',
  `user_id` int(5) NOT NULL DEFAULT '0' COMMENT '关联用户id',
  `status` int(2) NOT NULL DEFAULT '0' COMMENT '状态 0 禁用 1 启用',
  `addTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '入库时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=17 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='设备绑定房间信息表';

-- ----------------------------
-- Records of t_device_bind_info
-- ----------------------------
INSERT INTO `t_device_bind_info` VALUES ('10', 'dtl_divece_id', '1002', '15', '0', '2018-05-21 23:06:00');
INSERT INTO `t_device_bind_info` VALUES ('11', 'divece_mac', '501', '15', '0', '2018-05-23 21:41:50');
INSERT INTO `t_device_bind_info` VALUES ('12', '2368_Andriod', '1001', '24', '0', '2018-05-30 19:01:44');
INSERT INTO `t_device_bind_info` VALUES ('13', '753012345', '1001', '26', '0', '2018-05-30 20:59:45');
INSERT INTO `t_device_bind_info` VALUES ('14', '00090909098', '1001', '26', '0', '2018-06-14 07:35:05');
INSERT INTO `t_device_bind_info` VALUES ('15', '7530wwr0', '1000', '26', '0', '2018-05-31 00:05:04');
INSERT INTO `t_device_bind_info` VALUES ('16', '7530lily', '1006', '26', '0', '2018-05-31 00:05:50');

-- ----------------------------
-- Table structure for t_device_info
-- ----------------------------
DROP TABLE IF EXISTS `t_device_info`;
CREATE TABLE `t_device_info` (
  `id` int(4) NOT NULL AUTO_INCREMENT,
  `device_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '设备名称',
  `device_id` varchar(32) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '设备id，唯一',
  `status` int(2) NOT NULL DEFAULT '0' COMMENT '设备是否在线状态 1 在线 2 不在线',
  `barry` double NOT NULL DEFAULT '0' COMMENT '设备的电量',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户id',
  `gw_id` int(11) DEFAULT '0' COMMENT '网关id',
  `addTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=25 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='设备表信息';

-- ----------------------------
-- Records of t_device_info
-- ----------------------------
INSERT INTO `t_device_info` VALUES ('13', '电子表', 'divece_mac', '0', '0', '15', '8', '2018-05-22 21:27:10');
INSERT INTO `t_device_info` VALUES ('12', '门锁', 'dtl_divece_id', '0', '0', '15', '8', '2018-05-23 21:55:14');
INSERT INTO `t_device_info` VALUES ('14', '安卓系列', '2368_Andriod', '0', '0', '24', '9', '2018-05-30 19:01:44');
INSERT INTO `t_device_info` VALUES ('15', 'Room1', '753012345', '0', '0', '26', '10', '2018-05-30 20:59:45');
INSERT INTO `t_device_info` VALUES ('16', '测试', '00090909098', '0', '0', '26', '11', '2018-06-14 07:35:36');
INSERT INTO `t_device_info` VALUES ('23', 'wwr', '7530wwr0', '0', '0', '26', '11', '2018-05-31 00:05:04');
INSERT INTO `t_device_info` VALUES ('24', 'oioyoi', '7530lily', '0', '0', '26', '10', '2018-05-31 00:05:50');

-- ----------------------------
-- Table structure for t_device_open_info
-- ----------------------------
DROP TABLE IF EXISTS `t_device_open_info`;
CREATE TABLE `t_device_open_info` (
  `id` int(4) NOT NULL AUTO_INCREMENT,
  `device_id` varchar(32) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '设备id',
  `method_id` int(2) NOT NULL DEFAULT '1' COMMENT '开门方式 1 微信开门 2 滴卡开门 3 钥匙开门',
  `open_time` int(11) NOT NULL DEFAULT '0' COMMENT '开门时间戳',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='设备开门信息记录表';

-- ----------------------------
-- Records of t_device_open_info
-- ----------------------------
INSERT INTO `t_device_open_info` VALUES ('3', 'dtl_suo', '1', '1527402000');
INSERT INTO `t_device_open_info` VALUES ('4', 'dtl_divece_id02', '2', '1527402000');
INSERT INTO `t_device_open_info` VALUES ('5', 'dtl_divece_id', '1', '1527402000');
INSERT INTO `t_device_open_info` VALUES ('6', '7530wwr', '1', '1528935418');

-- ----------------------------
-- Table structure for t_gateway_info
-- ----------------------------
DROP TABLE IF EXISTS `t_gateway_info`;
CREATE TABLE `t_gateway_info` (
  `id` int(4) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT ' 网关名称',
  `gateway_id` varchar(32) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '网关设备id',
  `status` int(2) NOT NULL DEFAULT '0' COMMENT '设备是否在线  1 在线  0 不在线',
  `user_id` int(4) NOT NULL DEFAULT '0' COMMENT '用户关联 id',
  `addTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '数据入库时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=12 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='网关信息表';

-- ----------------------------
-- Records of t_gateway_info
-- ----------------------------
INSERT INTO `t_gateway_info` VALUES ('8', 'MAC', 'wang_001', '0', '15', '2018-05-22 21:27:10');
INSERT INTO `t_gateway_info` VALUES ('9', 'Android', 'dtl_way_Andriod', '0', '24', '2018-05-30 19:01:44');
INSERT INTO `t_gateway_info` VALUES ('10', 'tesst', '111111', '1', '26', '2018-06-12 17:47:14');
INSERT INTO `t_gateway_info` VALUES ('11', 'wega', '222222', '0', '26', '2018-06-12 17:23:54');

-- ----------------------------
-- Table structure for t_manger_pushsetting_info
-- ----------------------------
DROP TABLE IF EXISTS `t_manger_pushsetting_info`;
CREATE TABLE `t_manger_pushsetting_info` (
  `id` int(4) NOT NULL AUTO_INCREMENT,
  `user_id` int(4) NOT NULL DEFAULT '0' COMMENT '用户id',
  `url` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '数据推送接口设置',
  `token_url` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT 'token 获取地址',
  `appid` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT 'appid',
  `secret` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT 'key',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '数据入库',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='管理用户数据推送功能基本设置';

-- ----------------------------
-- Records of t_manger_pushsetting_info
-- ----------------------------
INSERT INTO `t_manger_pushsetting_info` VALUES ('1', '15', 'www.help.com', 'www.help.com?token=13555fsdf', 'b91efcfee87307d14a3ab21c77fe5ee1', '95365676dcde7d555e09abe451b30745', '1527318276');
INSERT INTO `t_manger_pushsetting_info` VALUES ('4', '24', 'http://hao123.com', 'http://hao123.com?u=123456', 'vzcvzvvzxvxbxcvbxcvbxcvb', 'bcxvbxcvbxcvbxcbvnbvnvncvncn', '1527904137');

-- ----------------------------
-- Table structure for t_open_method_info
-- ----------------------------
DROP TABLE IF EXISTS `t_open_method_info`;
CREATE TABLE `t_open_method_info` (
  `id` int(4) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '开门方式',
  `addTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='开门方式基本信息表';

-- ----------------------------
-- Records of t_open_method_info
-- ----------------------------
INSERT INTO `t_open_method_info` VALUES ('1', '微信开门', '2018-05-12 10:54:33');
INSERT INTO `t_open_method_info` VALUES ('2', '滴卡开门', '2018-05-12 10:54:40');
INSERT INTO `t_open_method_info` VALUES ('3', '钥匙开门', '2018-05-12 10:54:52');
INSERT INTO `t_open_method_info` VALUES ('4', '密码开门', '2018-06-14 06:51:07');

-- ----------------------------
-- Table structure for t_product_info
-- ----------------------------
DROP TABLE IF EXISTS `t_product_info`;
CREATE TABLE `t_product_info` (
  `id` int(4) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '类型 ',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '入库时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='产品类型';

-- ----------------------------
-- Records of t_product_info
-- ----------------------------
INSERT INTO `t_product_info` VALUES ('1', '网关', '2018-05-12 10:56:03');
INSERT INTO `t_product_info` VALUES ('2', '门锁模块', '2018-05-12 10:56:12');

-- ----------------------------
-- Table structure for t_staff_system
-- ----------------------------
DROP TABLE IF EXISTS `t_staff_system`;
CREATE TABLE `t_staff_system` (
  `ID` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `User_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
  `System_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '网站id',
  `Created_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `Updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '修改时间',
  PRIMARY KEY (`ID`),
  UNIQUE KEY `User_id` (`User_id`,`System_id`) USING BTREE
) ENGINE=MyISAM AUTO_INCREMENT=1343 DEFAULT CHARSET=utf8 COMMENT='用户权限表';

-- ----------------------------
-- Records of t_staff_system
-- ----------------------------

-- ----------------------------
-- Table structure for t_system_manage
-- ----------------------------
DROP TABLE IF EXISTS `t_system_manage`;
CREATE TABLE `t_system_manage` (
  `id` int(4) NOT NULL AUTO_INCREMENT,
  `parent_id` int(4) NOT NULL DEFAULT '0' COMMENT '父级id',
  `title` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '标题',
  `path` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '路径',
  `remark` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '描述',
  `sort` int(4) NOT NULL DEFAULT '0' COMMENT '排序',
  `Is_show` int(2) NOT NULL DEFAULT '0' COMMENT '显示 0 不显示 1显示',
  `addTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '入库时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='功能列表';

-- ----------------------------
-- Records of t_system_manage
-- ----------------------------

-- ----------------------------
-- Table structure for t_user_info
-- ----------------------------
DROP TABLE IF EXISTS `t_user_info`;
CREATE TABLE `t_user_info` (
  `id` int(4) NOT NULL AUTO_INCREMENT,
  `agent_id` int(20) DEFAULT '0' COMMENT '代理商的前缀ID',
  `user_account` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '登录名',
  `user_pwd` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '登陆密码',
  `parent_id` int(4) NOT NULL DEFAULT '0' COMMENT '添加用户的id',
  `user_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '用户名',
  `user_addr` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '用户地址',
  `user_phone` varchar(11) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '用户电话',
  `user_log` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '图片',
  `user_type` int(4) NOT NULL DEFAULT '0' COMMENT '用户类型 关联 t_user_type_info',
  `is_edit` int(2) NOT NULL DEFAULT '0' COMMENT '是否上转基本信息 0 未上传过 1已经上传过（上传过本人不能进行修改） ',
  `limit_status` int(2) NOT NULL DEFAULT '0' COMMENT '状态 0 可用 1 禁用',
  `appid` varchar(128) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '创建用户自动生成',
  `secret` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '创建用户自动生成',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '入库时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=27 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='用户基本信息';

-- ----------------------------
-- Records of t_user_info
-- ----------------------------
INSERT INTO `t_user_info` VALUES ('1', '8888', 'admin', 'f379eaf3c831b04de153469d1bec345e', '-1', 'admin', '深圳市南山区科研路', '15177957381', '/uploadfiles/2018/152553610198313/t01dfe76d16363bb0c2.jpg', '1', '1', '1', '', '', '1525594476');
INSERT INTO `t_user_info` VALUES ('18', '0', 'test', 'f379eaf3c831b04de153469d1bec345e', '17', '酒店管理员', '深圳市宝安区', '15177957381', '/uploadfiles/2018/152671515078499/t013658e41e8c191970.jpg', '3', '0', '1', '3ce7a015b2d5b31a8b5ba7fe677d3649', 'db43376f1b1aa35642eda17c560437e4', '1526715150');
INSERT INTO `t_user_info` VALUES ('15', '1234', 'gly01', 'e44d44aeb854d5578dded9b7cf1beeed', '17', '深圳市酒店管理员', '深圳市南山区科苑路', '13823528937', '', '3', '1', '1', '6fa0353c6b3bde4e27b44af3825e880c', '4d6bf737d7c369f91266ee6a12770f55', '1526483609');
INSERT INTO `t_user_info` VALUES ('17', '2487', 'admin', 'f379eaf3c831b04de153469d1bec345e', '1', '代理商测试', '深圳市宝安区沙井街道', '13823528937', '/uploadfiles/2018/152671447499285/t018c5c82327ff13b70.jpg', '2', '0', '1', 'ca4266ab8fcb02e2ea1bdb350de6d1f4', 'd093173074f7d37b8e04d9bc99897483', '1527510391');
INSERT INTO `t_user_info` VALUES ('19', '2368', 'yf', 'f379eaf3c831b04de153469d1bec345e', '1', 'yf', '深圳市南山区科苑北', '13823528937', '/uploadfiles/2018/152750998291988/t013658e41e8c191970.jpg', '2', '1', '1', '3c2506ae7e5f5ab0be7afb9cf99bc1e0', '562eeda26363a4e166d7ed79be2ccf27', '1527677350');
INSERT INTO `t_user_info` VALUES ('20', '6570', 'test002', 'f379eaf3c831b04de153469d1bec345e', '1', 'test002', '深圳市宝安区西乡', '15177957381', '/uploadfiles/2018/152751044721403/t018c5c82327ff13b70.jpg', '2', '0', '1', '4fcf088b1455e626a2668218d50895df', '97c2f00fc077d1a0b1e9c81127632655', '1527562856');
INSERT INTO `t_user_info` VALUES ('23', '7530', 'dls001', 'f379eaf3c831b04de153469d1bec345e', '1', 'test001', '深圳市南山区', '15177957381', '/uploadfiles/2018/152767744425004/t013658e41e8c191970.jpg', '2', '0', '1', '405e72bc499ec637f27dadb7e7d920b4', '8b9d0263b5e7253fba89180f86ae5f1a', '1527677444');
INSERT INTO `t_user_info` VALUES ('24', '0', 'gly01', 'f379eaf3c831b04de153469d1bec345e', '19', '管理员test', '深圳市南山区', '13823528930', '/uploadfiles/2018/152767795286645/t013658e41e8c191970.jpg', '3', '1', '1', '36e3bc24f1bad922ebab9c794069b8b0', '679a05e91edca36b45ae0fe1de654568', '1527678607');
INSERT INTO `t_user_info` VALUES ('25', '7530', 'jking', '0192023a7bbd73250516f069df18b500', '1', 'wenzhongjian', 'fasfw', '13723450181', '/uploadfiles/2018/152768431492460/微信图片_20180509224932.png', '2', '0', '1', '8dddeeffbf83d067fa4f385e33f15a3f', 'ba71399e517572dbba14a1360ea86900', '1527684314');
INSERT INTO `t_user_info` VALUES ('26', '0', 'hotel1', '0192023a7bbd73250516f069df18b500', '25', 'wenzhongjian', '深圳市南山区', '13723450181', '/uploadfiles/2018/152768450923990/微信图片_20180509224932.png', '3', '0', '1', '88236f6aa303fb212da966f5fe275b7c', 'd74270fbf70eaa517b8279a29b5db083', '1527684509');

-- ----------------------------
-- Table structure for t_user_type_info
-- ----------------------------
DROP TABLE IF EXISTS `t_user_type_info`;
CREATE TABLE `t_user_type_info` (
  `id` int(4) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '用户类型',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '入库时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='用户类型';

-- ----------------------------
-- Records of t_user_type_info
-- ----------------------------
INSERT INTO `t_user_type_info` VALUES ('1', '超级管理员', '1525354504');
INSERT INTO `t_user_type_info` VALUES ('2', '代理商', '1525354504');
INSERT INTO `t_user_type_info` VALUES ('3', '管理用户', '1525354504');
