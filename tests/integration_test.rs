use std::env;
use std::sync::Arc;

use sea_query::TableRenameStatement;
use snake::models::kline::{get_kline_table_name, KlineType};
use snake::pkg::log::logger::init_logger;
use snake::pkg::log::prelude::*;
use snake::pkg::sea_orm::prelude::*;
use snake::storage::mysql::kline;
use tokio;

#[tokio::test]
async fn test_crate_kline_db() {
    env::set_var("MYSQL_USER", "alex");
    let log_cfg = LogConfig::default();
    init_logger(&log_cfg);

    let mut orm_cfg = OrmConfig::default();
    orm_cfg.sqlx_logging = true;
    orm_cfg.sqlx_logging_level = "debug".to_string();
    orm_cfg.update_by_env();

    let conn = connect(&orm_cfg).await.unwrap();
    let conn = Arc::new(conn);

    let symbols = vec!["BTCUSDT", "ETHUSDT", "BTCUSD"];

    let kline_types = vec![
        KlineType::Min1,
        KlineType::Min5,
        KlineType::Min15,
        KlineType::Min30,
        KlineType::Hour1,
        KlineType::Hour4,
        KlineType::Day1,
        KlineType::Week1,
        KlineType::Month1,
    ];

    let mut table_names = Vec::new();

    symbols.iter().all(|&symbol| {
        kline_types.iter().all(|&kline_type| {
            table_names.push(get_kline_table_name(symbol, kline_type));
            true
        })
    });

    for name in table_names {
        kline::create_kline_table(&conn, &name).await.unwrap();
    }
}
