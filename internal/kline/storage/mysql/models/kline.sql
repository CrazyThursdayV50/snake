CREATE TABLE `kline` (
  `open_ts` bigint unsigned NOT NULL,
  `close_ts` bigint unsigned NOT NULL,
  `open` varchar(40) NOT NULL DEFAULT "0",
  `close` varchar(40) NOT NULL DEFAULT "0",
  `low` varchar(40) NOT NULL DEFAULT "0",
  `high` varchar(40) NOT NULL DEFAULT "0",
  `average` varchar(40) NOT NULL DEFAULT "0",
  `volume` varchar(40) NOT NULL DEFAULT "0",
  `amount` varchar(40) NOT NULL DEFAULT "0",
  `trade_count` int unsigned NOT NULL DEFAULT "0",
  `taker_buy_volume` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `taker_buy_amount` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`open_ts`),
  UNIQUE KEY `uk_kline_btcusdt_1m_close_ts` (`close_ts`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
