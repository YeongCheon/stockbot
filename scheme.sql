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
DROP TABLE IF EXISTS `user`;
CREATE TABLE IF NOT EXISTS `user`(
       `email` VARCHAR(30) PRIMARY KEY,
       `name` VARCHAR(30) NOT NULL
);
INSERT INTO user(`email`, `name`) VALUES(`kyc1682@gmail.com`, `YeongCheon`)

/*주식 등락기록*/
DROP TABLE IF EXISTS `stock_log`;
CREATE TABLE IF NOT EXISTS `stock_log`(
       `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
       `stock_code` CHAR(6) REFERENCES stock(code) NOT NULL,
       `stock_market` CHAR(6) REFERENCES stock(market) NOT NULL,
       `ask` NUMERIC(15,2) NOT NULL, /*매도가*/
       `bid` NUMERIC(15,2) NOT NULL, /*매수가*/
       `datetime` TIMESTAMP NOT NULL
);

/*주식 거래기록*/
DROP TABLE IF EXISTS `trade_log`;
CREATE TABLE IF NOT EXISTS `trade_log`(
       `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
       `user_email` VARCHAR(30) REFERENCES user(`email`) NOT NULL,
       `stock_code` CHAR(6) REFERENCES stock(code) NOT NULL,
       `stock_market` CHAR(6) REFERENCES stock(market) NOT NULL,
       `trade_type` CHAR(3) CHECK (`trade_type` IN('ask','bid')) NOT NULL,
       `price` NUMERIC(15,2) NOT NULL, /*거래체결 가격*/
       `trade_timestamp` timestamp DEFAULT CURRENT_TIMESTAMP
);


/*사용자의 주식 보유정보*/
DROP TABLE IF EXISTS `user_stock`;
CREATE TABLE IF NOT EXISTS `user_stock`(
       `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
       `user_email` VARCHAR(30) REFERENCES user(`email`) NOT NULL,
       `stock_code` CHAR(6) REFERENCES stock(code) NOT NULL
);
