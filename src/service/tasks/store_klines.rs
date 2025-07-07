use std::sync::Arc;
use std::time::Duration;

use super::db_kline_consumer::Consumer as DBKlineConsumer;
use crate::cexs::binance::prelude::*;
use crate::storage::mysql::models::kline_1m::{self, Kline, KlineInterval};
use sea_orm::*;
use tokio::time::sleep;

pub struct StoreKlineWorker {
    symbol: String,
    consumer: Arc<DBKlineConsumer>,
    db: Arc<DbConn>,
    client: Arc<MarketClient>,
}

impl StoreKlineWorker {
    pub fn new(
        symbol: String,
        db: Arc<DbConn>,
        client: Arc<MarketClient>,
        consumer: Arc<DBKlineConsumer>,
    ) -> Self {
        Self {
            symbol,
            db,
            client,
            consumer,
        }
    }

    async fn store_kline_from(&self, start_time: &mut i64, stop_time: Option<i64>) -> bool {
        let klines = self
            .client
            .get_klines_from(
                &self.symbol,
                *start_time,
                stop_time,
                kline_1m::KlineInterval::Min1,
                false,
            )
            .await;

        if klines.len() > 0 {
            let next_start_time = klines.clone().last().unwrap().open_ts.clone().unwrap();
            Kline::insert_many(klines)
                .exec(self.db.as_ref())
                .await
                .unwrap();

            *start_time = next_start_time;
            true
        } else {
            false
        }
    }

    async fn store_kline_to(&self, end_time: &mut i64) -> bool {
        let klines = self
            .client
            .get_klines_to(
                &self.symbol,
                *end_time,
                None,
                kline_1m::KlineInterval::Min1,
                false,
            )
            .await;

        if klines.len() > 0 {
            let next_end_time = klines.clone().first().unwrap().open_ts.clone().unwrap();
            Kline::insert_many(klines)
                .exec(self.db.as_ref())
                .await
                .unwrap();

            *end_time = next_end_time;
            true
        } else {
            false
        }
    }

    pub async fn run_store_klines(&self, stop_time: i64) {
        // find latest kline
        let mut start_time = 0i64;
        let mut end_time = 0i64;

        if let Ok(Some(latest_kline)) = Kline::find()
            .order_by(kline_1m::Column::OpenTs, Order::Desc)
            .one(self.db.as_ref())
            .await
        {
            start_time = latest_kline.open_ts;
            end_time = latest_kline.open_ts;
            // 如果存在最后一条 k 线，那么补齐 start_time 到 stop_time 之间的 k 线
            let mut next = true;
            while next {
                if start_time >= stop_time {
                    break;
                }

                next = self
                    .store_kline_from(&mut start_time, Some(stop_time))
                    .await;
                _ = sleep(Duration::from_secs(2)).await;
            }

            if let Ok(Some(earliest_kline)) = Kline::find()
                .order_by(kline_1m::Column::OpenTs, Order::Asc)
                .one(self.db.as_ref())
                .await
            {
                end_time = earliest_kline.open_ts
            };

            // 从 end_time 开始，往前拿到所有 k 线
            next = true;
            while next {
                next = self.store_kline_to(&mut end_time).await;
                _ = sleep(Duration::from_secs(2)).await;
            }
        }
        // Implementation of store_klines function
    }
}
