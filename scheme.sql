DROP TABLE IF EXISTS `stockbot_kospi`;
CREATE TABLE IF NOT EXISTS `stockbot_kospi`(
       `kospi_code` CHAR(6) PRIMARY KEY,
       `kospi_name` VARCHAR(20) NOT NULL,
       `kospi_category_code` CHAR(6) NOT NULL,
       `kospi_category` VARCHAR(50) NOT NULL,
       `kospi_stock_total` BIGINT UNSIGNED DEFAULT 0,
       `kospi_capital` NUMERIC(15,2),
       `kospi_facevalue` NUMERIC(15,2),
       `kospi_currency` CHAR(10),
       `kospi_tel` CHAR(20),
       `kospi_address` VARCHAR(50),
       `kospi_totalcount` INT UNSIGNED
);

DROP TABLE IF EXISTS `stockbo_board`;
CREATE TABLE IF NOT EXISTS `stockbot_board`(
       `board_id` CHAR(9) PRIMARY KEY, /*for yahoo finance.(ex : 005930.KS)*/
       `board_name` VARCHAR(20)  DEFAULT NULL
);

DROP TABLE IF EXISTS `stockbot_kospi_history`;
CREATE TABLE IF NOT EXISTS `stockbot_kospi_history`(
       `history_id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
       `history_code` CHAR(9) REFERENCES stockbot_kospi(kospi_code),
       `history_ask` NUMERIC(15,2) NOT NULL,
       `history_bid` NUMERIC(15,2) NOT NULL,
       `history_datetime` TIMESTAMP
);
