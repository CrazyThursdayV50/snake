use crate::{cexs::binance::prelude::*, storage::mysql::models::kline_1m::KlineInterval};
use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Config {
    pub cexs: Cex,
    pub klines: Vec<Kline>,
}

#[derive(Debug, Deserialize)]
pub struct Kline {
    pub symbol: String,
    pub intervals: Vec<KlineInterval>,
}

#[derive(Debug, Deserialize)]
pub struct Cex {
    pub binance: BinanceConfig,
}
