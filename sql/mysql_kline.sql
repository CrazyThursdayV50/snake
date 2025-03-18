-- MySQL K线表示例
CREATE TABLE IF NOT EXISTS `mysql_kline` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `symbol` VARCHAR(20) NOT NULL COMMENT '交易对',
  `period` VARCHAR(10) NOT NULL COMMENT '周期',
  `open_ts` BIGINT NOT NULL COMMENT '开盘时间戳',
  `close_ts` BIGINT NOT NULL COMMENT '收盘时间戳',
  `open` DOUBLE NOT NULL COMMENT '开盘价',
  `high` DOUBLE NOT NULL COMMENT '最高价',
  `low` DOUBLE NOT NULL COMMENT '最低价',
  `close` DOUBLE NOT NULL COMMENT '收盘价',
  `volume` DOUBLE NOT NULL COMMENT '成交量',
  `amount` DOUBLE NOT NULL COMMENT '成交额',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_mysql_kline_symbol_period_close_ts` (`symbol`, `period`, `close_ts`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='MySQL K线表示例'; 