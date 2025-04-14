CREATE DATABASE IF NOT EXISTS responsible_api;
drop user 'responsible_api_user'@'%';
CREATE USER 'responsible_api_user'@'%' IDENTIFIED BY 'responsible_api_pass';
GRANT ALL PRIVILEGES ON responsible_api.* TO 'responsible_api_user'@'%';
FLUSH PRIVILEGES;
USE responsible_api;

-- Create syntax for TABLE 'responsible_api_users'
CREATE TABLE
  `responsible_api_users` (
    `uid` int unsigned NOT NULL AUTO_INCREMENT,
    `account_id` bigint NOT NULL DEFAULT '0',
    `name` varchar(60) NOT NULL DEFAULT '',
    `mail` varchar(254) DEFAULT '',
    `created` int NOT NULL DEFAULT '0',
    `access` int NOT NULL DEFAULT '0',
    `status` tinyint NOT NULL DEFAULT '0',
    `secret` varchar(32) NOT NULL DEFAULT '',
    `apikey` varchar(64) DEFAULT '',
    `refresh_token` varchar(128) DEFAULT '',
    PRIMARY KEY (`uid`),
    UNIQUE KEY `name` (`name`),
    KEY `access` (`access`),
    KEY `created` (`created`),
    KEY `mail` (`mail`),
    KEY `account_id` (`account_id`)
  ) ENGINE = InnoDB;

-- Create syntax for TABLE 'responsible_token_bucket'
CREATE TABLE
  `responsible_token_bucket` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `bucket` varchar(128) NOT NULL DEFAULT '',
    `account_id` bigint DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY `Account ID Constraint` (`account_id`),
    CONSTRAINT `Account ID Constraint` FOREIGN KEY (`account_id`) REFERENCES `responsible_api_users` (`account_id`)
  ) ENGINE = InnoDB;