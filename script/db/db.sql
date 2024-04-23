/*
Navicat MySQL Data Transfer

Source Server         : local
Source Server Version : 80300
Source Host           : localhost:3306
Source Database       : go-backend

Target Server Type    : MYSQL
Target Server Version : 80300
File Encoding         : 65001

Date: 2024-04-23 12:06:36
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for interaction
-- ----------------------------
DROP TABLE IF EXISTS `interaction`;
CREATE TABLE `interaction` (
                               `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                               `created_at` datetime(3) DEFAULT NULL,
                               `updated_at` datetime(3) DEFAULT NULL,
                               `deleted_at` datetime(3) DEFAULT NULL,
                               `biz_id` bigint DEFAULT NULL,
                               `biz` varchar(255) DEFAULT NULL,
                               `read_cnt` bigint DEFAULT NULL,
                               `like_cnt` bigint DEFAULT NULL,
                               `collect_cnt` bigint DEFAULT NULL,
                               PRIMARY KEY (`id`),
                               UNIQUE KEY `idx_interaction_bizID_biz` (`biz_id`),
                               UNIQUE KEY `idx_bizID_biz` (`biz`),
                               KEY `idx_interaction_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of interaction
-- ----------------------------

-- ----------------------------
-- Table structure for post
-- ----------------------------
DROP TABLE IF EXISTS `post`;
CREATE TABLE `post` (
                        `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                        `created_at` datetime(3) DEFAULT NULL,
                        `updated_at` datetime(3) DEFAULT NULL,
                        `deleted_at` datetime(3) DEFAULT NULL,
                        `post_id` bigint NOT NULL,
                        `title` longtext,
                        `abstract` longtext,
                        `content` longtext,
                        `author_id` bigint DEFAULT NULL,
                        `status` tinyint unsigned DEFAULT NULL,
                        PRIMARY KEY (`id`,`post_id`),
                        KEY `idx_post_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of post
-- ----------------------------
INSERT INTO `post` VALUES ('1', '2024-04-23 12:01:29.402', '2024-04-23 12:01:29.402', null, '173282172709376000', '1111111111', '1111111111', '1111111111', '173282047035445248', '1');
INSERT INTO `post` VALUES ('2', '2024-04-23 12:01:33.190', '2024-04-23 12:01:33.190', null, '173282188601593856', '111111111', '111111111', '111111111', '173282047035445248', '1');
INSERT INTO `post` VALUES ('3', '2024-04-23 12:01:36.125', '2024-04-23 12:01:36.125', null, '173282200907681792', '11111111', '11111111', '11111111', '173282047035445248', '1');
INSERT INTO `post` VALUES ('4', '2024-04-23 12:01:40.093', '2024-04-23 12:01:40.093', null, '173282217554874368', '1111111', '1111111', '1111111', '173282047035445248', '1');
INSERT INTO `post` VALUES ('5', '2024-04-23 12:01:43.663', '2024-04-23 12:01:43.663', null, '173282232520151040', '111111', '111111', '111111', '173282047035445248', '1');
INSERT INTO `post` VALUES ('6', '2024-04-23 12:01:46.577', '2024-04-23 12:01:46.577', null, '173282244750741504', '11111', '11111', '11111', '173282047035445248', '1');
INSERT INTO `post` VALUES ('7', '2024-04-23 12:01:50.427', '2024-04-23 12:01:50.427', null, '173282260894617600', '1111', '1111', '1111', '173282047035445248', '1');
INSERT INTO `post` VALUES ('8', '2024-04-23 12:01:53.525', '2024-04-23 12:01:53.525', null, '173282273888571392', '111', '111', '111', '173282047035445248', '1');
INSERT INTO `post` VALUES ('9', '2024-04-23 12:01:56.615', '2024-04-23 12:01:56.615', null, '173282286848970752', '11', '11', '11', '173282047035445248', '1');
INSERT INTO `post` VALUES ('10', '2024-04-23 12:01:59.631', '2024-04-23 12:01:59.631', null, '173282299498991616', '1', '1', '1', '173282047035445248', '1');

-- ----------------------------
-- Table structure for task
-- ----------------------------
DROP TABLE IF EXISTS `task`;
CREATE TABLE `task` (
                        `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                        `created_at` datetime(3) DEFAULT NULL,
                        `updated_at` datetime(3) DEFAULT NULL,
                        `deleted_at` datetime(3) DEFAULT NULL,
                        `task_id` bigint DEFAULT NULL,
                        `title` longtext,
                        `user_id` bigint DEFAULT NULL,
                        PRIMARY KEY (`id`),
                        KEY `idx_task_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of task
-- ----------------------------

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
                        `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                        `created_at` datetime(3) DEFAULT NULL,
                        `updated_at` datetime(3) DEFAULT NULL,
                        `deleted_at` datetime(3) DEFAULT NULL,
                        `user_id` bigint NOT NULL,
                        `user_name` varchar(191) DEFAULT NULL,
                        `email` varchar(191) DEFAULT NULL,
                        `password` varchar(256) DEFAULT NULL,
                        PRIMARY KEY (`id`,`user_id`),
                        UNIQUE KEY `uni_user_user_name` (`user_name`),
                        UNIQUE KEY `uni_user_email` (`email`),
                        KEY `idx_user_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES ('1', '2024-04-23 12:00:48.817', '2024-04-23 12:00:48.817', null, '173282002483548160', 'admin', 'test@qq.com', '$2a$10$Nv.c3kA/Y/1myZ9qFRz3LOiQN95QSUhV72bzrvdcv/LBiFNwBaZKq');
INSERT INTO `user` VALUES ('2', '2024-04-23 12:00:59.440', '2024-04-23 12:00:59.440', null, '173282047035445248', 'root', '1227891082@qq.com', '$2a$10$h94fcM7EQXRQcyeuN6PhW.r0a2sblpdmm/C/aIJ1cls4pLJlszlmK');

-- ----------------------------
-- Table structure for user_collect
-- ----------------------------
DROP TABLE IF EXISTS `user_collect`;
CREATE TABLE `user_collect` (
                                `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                                `created_at` datetime(3) DEFAULT NULL,
                                `updated_at` datetime(3) DEFAULT NULL,
                                `deleted_at` datetime(3) DEFAULT NULL,
                                `user_id` bigint DEFAULT NULL,
                                `biz_id` bigint DEFAULT NULL,
                                `biz` varchar(255) DEFAULT NULL,
                                `collection_id` bigint DEFAULT NULL,
                                `status` tinyint(1) DEFAULT NULL,
                                PRIMARY KEY (`id`),
                                UNIQUE KEY `idx_userCollect_userID_bizID_biz` (`user_id`,`biz_id`,`biz`),
                                KEY `idx_user_collect_deleted_at` (`deleted_at`),
                                KEY `idx_user_collect_collection_id` (`collection_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of user_collect
-- ----------------------------

-- ----------------------------
-- Table structure for user_like
-- ----------------------------
DROP TABLE IF EXISTS `user_like`;
CREATE TABLE `user_like` (
                             `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                             `created_at` datetime(3) DEFAULT NULL,
                             `updated_at` datetime(3) DEFAULT NULL,
                             `deleted_at` datetime(3) DEFAULT NULL,
                             `user_id` bigint DEFAULT NULL,
                             `biz_id` bigint DEFAULT NULL,
                             `biz` varchar(255) DEFAULT NULL,
                             `status` tinyint(1) DEFAULT NULL,
                             PRIMARY KEY (`id`),
                             UNIQUE KEY `idx_userLike_userID_bizID_biz` (`user_id`,`biz_id`,`biz`),
                             KEY `idx_user_like_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of user_like
-- ----------------------------
