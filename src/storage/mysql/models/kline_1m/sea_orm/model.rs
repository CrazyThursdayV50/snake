use super::relation::Relation;
use binance::model::{KlineEvent, KlineSummary};
use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, Eq, DeriveEntityModel, Serialize, Deserialize)]
#[sea_orm(table_name = "kline_1m")]
pub struct Model {
    #[sea_orm(primary_key)]
    pub id: u64,
    pub symbol: String,
    pub open_ts: i64,
    pub close_ts: i64,
    pub open: String,
    pub close: String,
    pub low: String,
    pub high: String,
    pub average: String,
    pub volume: String,
    pub amount: String,
    pub trade_count: i64,
    pub taker_buy_volume: String,
    pub taker_buy_amount: String,
    pub created_at: DateTimeWithTimeZone,
    pub updated_at: DateTimeWithTimeZone,
}
