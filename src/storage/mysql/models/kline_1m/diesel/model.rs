use super::schema::klines;
use chrono::NaiveDateTime;
use diesel::prelude::*;

#[derive(Queryable, Selectable, Insertable)]
#[diesel(check_for_backend(diesel::mysql::Mysql))]
#[diesel(table_name = klines)] // 假设你的表名为 posts
pub struct Model {
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
    pub trade_count: i32,
    pub taker_buy_volume: String,
    pub taker_buy_amount: String,
    pub created_at: NaiveDateTime,
    pub updated_at: NaiveDateTime,
}
