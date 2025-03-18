use sea_orm::DbConn;

use crate::models::kline::get_kline_table_name;
use crate::service::config::Config;
use crate::storage::mysql::kline::create_kline_table;

pub async fn run(db: &DbConn, cfg: &Config) {
    for kline in cfg.klines.as_slice() {
        let table_name = get_kline_table_name(&kline.symbol, kline.interval);
        _ = create_kline_table(db, &table_name).await;
    }
}
