use std::sync::Arc;

use sea_orm::DbConn;

use crate::cexs::binance::prelude::MarketClient;
use crate::storage::mysql::models::kline_1m::KlineInterval;

use super::config::Config;
use super::tasks::*;

pub async fn init(db: Arc<DbConn>, client: Arc<MarketClient>, cfg: &Config) {
    let consumer = Arc::new(DBKlineConsumer::new(db));

    cfg.klines.iter().all(|kline| {
        let worker = StoreKlineWorker::new(
            kline.symbol.clone(),
            db.clone(),
            client.clone(),
            consumer.clone(),
        );

        client.subscribe_kline(&kline.symbol, KlineInterval::Min1);
        worker.run_store_klines(10);
        todo!("run worker");
        true
    });
}
