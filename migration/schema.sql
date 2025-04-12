CREATE DATABASE IF NOT EXISTS responsible_api;
drop user 'responsible_api_user'@'%';
CREATE USER 'responsible_api_user'@'%' IDENTIFIED BY 'responsible_api_pass';
GRANT ALL PRIVILEGES ON responsible_api_db.* TO 'responsible_api_user'@'%';
FLUSH PRIVILEGES;
USE responsible_api_db;