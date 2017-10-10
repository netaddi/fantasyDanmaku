CREATE DATABASE danmaku CHARACTER SET utf8mb4;

use danmaku;

CREATE TABLE users
(
  id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY ,
  username VARCHAR(128),
  password VARCHAR(128),
  nickname VARCHAR(128),
  reg_code VARCHAR(128) NOT NULL ,
  permission INTEGER,
  enrolled BOOLEAN NOT NULL
);

SELECT * FROM users;

INSERT INTO users (username, password, nickname, reg_code, permission, enrolled)
  VALUES ('root', 'testpass', 'root', '0', 1, TRUE );
