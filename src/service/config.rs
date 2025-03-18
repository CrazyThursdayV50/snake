use crate::models::kline::KlineType;
use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Config {
    pub klines: Vec<Kline>,
}

#[derive(Debug, Deserialize)]
pub struct Kline {
    pub symbol: String,
    pub interval: KlineType,
}
