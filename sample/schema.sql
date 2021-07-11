-- This file generated by protoc-ddl.
-- Generated by MySQL Generator (v0.1)

SET foreign_key_checks=0;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
	`id` INTEGER NOT NULL AUTO_INCREMENT,
	`age` INTEGER NOT NULL DEFAULT 20,
	`name` VARCHAR(255) NOT NULL,
	`title` VARCHAR(255) NOT NULL,
	`last_name` VARCHAR(255) NOT NULL,
	`created_at` DATETIME NOT NULL,
	UNIQUE `idx_title` (`title`),
	PRIMARY KEY(`id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `blog`;
CREATE TABLE `blog` (
	`id` BIGINT NOT NULL AUTO_INCREMENT,
	`user_id` INTEGER NOT NULL,
	`title` VARCHAR(100) NOT NULL,
	`body` TEXT NOT NULL,
	`category_id` INTEGER NULL,
	`attach` LONGBLOB NOT NULL,
	`editor_id` INTEGER NOT NULL,
	`sign` VARBINARY(20) NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NULL,
	UNIQUE `idx_user_id_and_id` (`user_id`, `id`),
	INDEX `idx_user_id_category_id` (`user_id`, `category_id`),
	UNIQUE `idx_user_id_title` (`user_id`, `title`),
	PRIMARY KEY(`id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `comment_image`;
CREATE TABLE `comment_image` (
	`comment_blog_id` BIGINT NOT NULL,
	`comment_user_id` INTEGER NOT NULL,
	`like_id` BIGINT UNSIGNED NOT NULL,
	PRIMARY KEY(`comment_blog_id`,`comment_user_id`,`like_id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
	`blog_id` BIGINT NOT NULL,
	`user_id` INTEGER NOT NULL,
	UNIQUE `idx_user_id` (`user_id`),
	PRIMARY KEY(`blog_id`,`user_id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `reply`;
CREATE TABLE `reply` (
	`id` INTEGER NOT NULL AUTO_INCREMENT,
	`comment_blog_id` BIGINT NULL,
	`comment_user_id` INTEGER NULL,
	`body` VARCHAR(255) NOT NULL,
	PRIMARY KEY(`id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `like`;
CREATE TABLE `like` (
	`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
	`user_id` INTEGER NOT NULL,
	`blog_id` BIGINT NOT NULL,
	PRIMARY KEY(`id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `post_image`;
CREATE TABLE `post_image` (
	`id` INTEGER NOT NULL,
	`url` VARCHAR(255) NOT NULL,
	PRIMARY KEY(`id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `task`;
CREATE TABLE `task` (
	`id` INTEGER NOT NULL AUTO_INCREMENT,
	`image_id` INTEGER NOT NULL,
	PRIMARY KEY(`id`)
) Engine=InnoDB;

SET foreign_key_checks=1;
