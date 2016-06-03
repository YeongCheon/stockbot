/*주식 종목정보*/
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

/*사용자 정보*/
DROP TABLE IF EXISTS `member`;
CREATE TABLE IF NOT EXISTS `member`(
       `email` VARCHAR(30) PRIMARY KEY,
       `name` VARCHAR(30) NOT NULL
);
INSERT INTO member(`email`, `name`) VALUES('kyc1682@gmail.com', 'YeongCheon');

/*주식 등락기록*/
DROP TABLE IF EXISTS `stock_log`;
CREATE TABLE IF NOT EXISTS `stock_log`(
       `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
       `stock_code` CHAR(6) NOT NULL,
       `stock_market` CHAR(6) NOT NULL,
       `ask` NUMERIC(15,2) NOT NULL, /*매도가*/
       `bid` NUMERIC(15,2) NOT NULL, /*매수가*/
       `datetime` TIMESTAMP NOT NULL,
       CONSTRAINT fk_stock_log FOREIGN KEY (`stock_code`,`stock_market`) REFERENCES `stock`(`code`, `market`) ON DELETE CASCADE ON UPDATE CASCADE
);


/*주식 거래기록*/
DROP TABLE IF EXISTS `trade_log`;
CREATE TABLE IF NOT EXISTS `trade_log`(
       `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
       `member_email` VARCHAR(30) NOT NULL,
       `stock_code` CHAR(6) NOT NULL,
       `stock_market` CHAR(6) NOT NULL,
       `trade_type` CHAR(3) NOT NULL CHECK (`trade_type` IN ('ask','bid')),
       `price` NUMERIC(15,2) NOT NULL, /*거래체결 가격*/
       `trade_timestamp` timestamp DEFAULT CURRENT_TIMESTAMP,
       CONSTRAINT fk_trade_log_email FOREIGN KEY (`member_email`) REFERENCES `member`(`email`) ON DELETE CASCADE ON UPDATE CASCADE,
       CONSTRAINT fk_trade_log_code FOREIGN KEY (`stock_code`, `stock_market`) REFERENCES `stock`(`code`, `market`) ON DELETE CASCADE ON UPDATE CASCADE
);


/*사용자의 주식 보유정보*/
DROP TABLE IF EXISTS `member_stock`;
CREATE TABLE IF NOT EXISTS `member_stock`(
       `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
       `member_email` VARCHAR(30) NOT NULL,
       `stock_code` CHAR(6) NOT NULL,
       `stock_market` CHAR(6) NOT NULL,
       CONSTRAINT fk_member FOREIGN KEY (`stock_code`, `stock_market`) REFERENCES `stock`(`code`, `market`) ON DELETE CASCADE ON UPDATE CASCADE
);
