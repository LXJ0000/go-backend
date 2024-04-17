CREATE DATABASE `go-backend`;
/*
 Navicat Premium Data Transfer

 Source Server         : 3306
 Source Server Type    : MySQL
 Source Server Version : 80027
 Source Host           : localhost:3306
 Source Schema         : go-backend

 Target Server Type    : MySQL
 Target Server Version : 80027
 File Encoding         : 65001

 Date: 13/04/2024 17:46:26
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for interaction
-- ----------------------------
DROP TABLE IF EXISTS `interaction`;
CREATE TABLE `interaction`  (
                                `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
                                `created_at` datetime(3) NULL DEFAULT NULL,
                                `updated_at` datetime(3) NULL DEFAULT NULL,
                                `deleted_at` datetime(3) NULL DEFAULT NULL,
                                `biz_id` bigint NULL DEFAULT NULL,
                                `biz` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
                                `read_cnt` bigint NULL DEFAULT NULL,
                                PRIMARY KEY (`id`) USING BTREE,
                                INDEX `idx_interaction_deleted_at`(`deleted_at` ASC) USING BTREE,
                                UNIQUE INDEX `idx_bizID_biz`(`biz_id` ASC, `biz` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of interaction
-- ----------------------------

-- ----------------------------
-- Table structure for post
-- ----------------------------
DROP TABLE IF EXISTS `post`;
CREATE TABLE `post`  (
                         `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
                         `created_at` datetime(3) NULL DEFAULT NULL,
                         `updated_at` datetime(3) NULL DEFAULT NULL,
                         `deleted_at` datetime(3) NULL DEFAULT NULL,
                         `post_id` bigint NOT NULL,
                         `title` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
                         `abstract` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
                         `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
                         `author_id` bigint NULL DEFAULT NULL,
                         `status` tinyint UNSIGNED NULL DEFAULT NULL,
                         PRIMARY KEY (`id`, `post_id`) USING BTREE,
                         INDEX `idx_post_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of post
-- ----------------------------
INSERT INTO `post` VALUES (1, '2024-04-13 17:44:45.626', '2024-04-13 17:44:45.626', NULL, 169744680886472704, '1111111', '1111111', '1111111', 169744592554430464, 1);
INSERT INTO `post` VALUES (2, '2024-04-13 17:45:41.497', '2024-04-13 17:45:41.497', NULL, 169744915230625792, '111111', '111111', '111111', 169744592554430464, 1);
INSERT INTO `post` VALUES (3, '2024-04-13 17:45:44.248', '2024-04-13 17:45:44.248', NULL, 169744926764961792, '11111', '11111', '11111', 169744592554430464, 1);
INSERT INTO `post` VALUES (4, '2024-04-13 17:45:46.532', '2024-04-13 17:45:46.532', NULL, 169744936348946432, '1111', '1111', '1111', 169744592554430464, 1);
INSERT INTO `post` VALUES (5, '2024-04-13 17:45:48.634', '2024-04-13 17:45:48.634', NULL, 169744945161179136, '111', '111', '111', 169744592554430464, 1);
INSERT INTO `post` VALUES (6, '2024-04-13 17:45:50.857', '2024-04-13 17:45:50.857', NULL, 169744954489311232, '11', '11', '11', 169744592554430464, 1);
INSERT INTO `post` VALUES (7, '2024-04-13 17:45:52.925', '2024-04-13 17:45:52.925', NULL, 169744963158937600, '1', '1', '1', 169744592554430464, 1);

-- ----------------------------
-- Table structure for task
-- ----------------------------
DROP TABLE IF EXISTS `task`;
CREATE TABLE `task`  (
                         `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
                         `created_at` datetime(3) NULL DEFAULT NULL,
                         `updated_at` datetime(3) NULL DEFAULT NULL,
                         `deleted_at` datetime(3) NULL DEFAULT NULL,
                         `task_id` bigint NULL DEFAULT NULL,
                         `title` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
                         `user_id` bigint NULL DEFAULT NULL,
                         PRIMARY KEY (`id`) USING BTREE,
                         INDEX `idx_task_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of task
-- ----------------------------

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
                         `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
                         `created_at` datetime(3) NULL DEFAULT NULL,
                         `updated_at` datetime(3) NULL DEFAULT NULL,
                         `deleted_at` datetime(3) NULL DEFAULT NULL,
                         `user_id` bigint NOT NULL,
                         `user_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
                         `email` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
                         `password` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
                         PRIMARY KEY (`id`, `user_id`) USING BTREE,
                         UNIQUE INDEX `uni_user_user_name`(`user_name` ASC) USING BTREE,
                         UNIQUE INDEX `uni_user_email`(`email` ASC) USING BTREE,
                         INDEX `idx_user_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (1, '2024-04-13 17:44:24.568', '2024-04-13 17:44:24.568', NULL, 169744592554430464, 'root', '1227891082@qq.com', '$2a$10$dOquVpWNAKxajgRRyniVLeywcPSSP3CMGilI7hcZgF7QcJdE.1glG');

SET FOREIGN_KEY_CHECKS = 1;
