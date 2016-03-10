DROP TABLE IF EXISTS `stock`;
CREATE TABLE IF NOT EXISTS `stock`(
       `code` CHAR(6),
       `market` CHAR(6),
       `name` VARCHAR(20) NOT NULL,
       `category_code` CHAR(6) NOT NULL,
       `category` VARCHAR(50) NOT NULL,
       `total` BIGINT UNSIGNED DEFAULT 0,
       `capital` NUMERIC(15,2),
       `facevalue` NUMERIC(15,2),
       `currency` CHAR(10),
       `tel` CHAR(20),
       `address` VARCHAR(50),
       `totalcount` INT UNSIGNED,
       PRIMARY KEY(code, market)
);

DROP TABLE IF EXISTS `user`;
CREATE TABLE IF NOT EXISTS `user`(
       `email` VARCHAR(30) PRIMARY KEY,
       `name` VARCHAR(30) NOT NULL
);

DROP TABLE IF EXISTS `stock_log`;
CREATE TABLE IF NOT EXISTS `stock_log`(
       `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
       `stock_code` CHAR(6) REFERENCES stock(code),
       `stock_market` CHAR(6) REFERENCES stock(market),
       `ask` NUMERIC(15,2) NOT NULL,
       `bid` NUMERIC(15,2) NOT NULL,
       `datetime` TIMESTAMP
);

DROP TABLE IF EXISTS `trade_log`;
CREATE TABLE IF NOT EXISTS `trade_log`(
       `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
       `user_email` VARCHAR(30) REFERENCES user(`email`),
       `stock_code` CHAR(6) REFERENCES stock(code),
       `stock_market` CHAR(6) REFERENCES stock(market),
       `trade_type` CHAR(3) CHECK (`trade_type` IN('ask','bid')),
       `price` NUMERIC(15,2) NOT NULL
);
