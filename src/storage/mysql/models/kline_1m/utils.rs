use super::KlineInterval;

const min_kline_limit: i64 = 3;

pub fn get_end_time_by_start_time(start_time: i64, kline_type: KlineInterval) -> i64 {
    start_time + (kline_type.duration().num_milliseconds() * (min_kline_limit - 1)) as i64
}

pub fn next_time(time: i64, kline_type: KlineInterval) -> i64 {
    time + kline_type.duration().num_milliseconds() as i64
}

pub fn last_time(time: i64, kline_type: KlineInterval) -> i64 {
    time - kline_type.duration().num_milliseconds() as i64
}

pub fn get_start_time_by_end_time(end_time: i64, kline_type: KlineInterval) -> i64 {
    end_time - (kline_type.duration().num_milliseconds() * (min_kline_limit - 1)) as i64
}

// pub fn stream_kline_from_str(stream_kline: &str) -> KlineEvent {
//     serde_json::from_str(stream_kline).unwrap()
// }

// #[cfg(test)]
// mod tests {
//     use binance::model::KlineSummary;
//     use binance_spot_connector_rust::http::request::Request;
//     use serde_json::Value;

//     use super::super::super::config::Config;
//     use super::super::client;
//     use super::*;
//     use crate::pkg::log::prelude::*;

//     #[tokio::test]
//     async fn test_kline() {
//         use std::env;

//         let log_cfg = LogConfig::default();
//         init_logger(&log_cfg);

//         let mut cfg = Config::default();
//         let api_key = env::var("BN_APIKEY").unwrap();
//         let secret = env::var("BN_SECRET").unwrap();
//         cfg.api_key = api_key;
//         cfg.secret_key = secret;
//         cfg.is_testnet = false;
//         let client = client::Client::new(&cfg);
//         let timestamp = 1737189420000;
//         let result = client
//             .get_klines_from("BTCUSDT", timestamp, None, KlineType::Min1, true)
//             .await;
//         // let value = to_value(result).unwrap();
//         // log::info!("value: {:?}", value);
//         // let klines: Vec<KlineSummary> = from_value(value).unwrap();
//         // log::info!("kline: {:?}", klines);
//         log::info!("kline: {:?}", result);
//     }
// }
