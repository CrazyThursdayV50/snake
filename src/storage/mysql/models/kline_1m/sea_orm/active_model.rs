use super::model::ActiveModel;
use binance::model::{KlineEvent, KlineSummary};
use rust_decimal::Decimal;
use sea_orm::prelude::*;
use sea_orm::*;
use std::str::FromStr;

impl ActiveModelBehavior for ActiveModel {}

use serde_json;
use serde_json::Value;

pub struct ActiveModelArray(Vec<ActiveModel>);

impl<T> From<T> for ActiveModelArray
where
    T: Into<String>,
{
    fn from(value: T) -> Self {
        let values: Vec<Vec<Value>> = serde_json::from_str(value.into().as_str()).unwrap();
        let klines: Vec<ActiveModel> = values
            .iter()
            .map(|v| KlineSummary::try_from(v).unwrap().into())
            .collect();

        ActiveModelArray(klines)
    }
}

impl ActiveModel {
    pub fn from_event<T>(value: T) -> Self
    where
        T: Into<String>,
    {
        let event: KlineEvent = serde_json::from_str(value.into().as_str()).unwrap();
        event.into()
    }
}

impl Into<Vec<ActiveModel>> for ActiveModelArray {
    fn into(self) -> Vec<ActiveModel> {
        self.0
    }
}

impl From<KlineEvent> for ActiveModel {
    fn from(event: KlineEvent) -> Self {
        let amount = Decimal::from_str(&event.kline.quote_asset_volume)
            .inspect_err(|e| log::error!("get amount from kline event: {:?}", e))
            .unwrap();

        let volume = Decimal::from_str(&event.kline.volume)
            .inspect_err(|e| log::error!("get volume from kline event: {:?}", e))
            .unwrap();

        let mut kline_model = Self::new();
        kline_model.symbol = Set(event.symbol);
        kline_model.open = Set(event.kline.open);
        kline_model.close = Set(event.kline.close);
        kline_model.high = Set(event.kline.high);
        kline_model.low = Set(event.kline.low);
        kline_model.volume = Set(event.kline.volume);
        kline_model.amount = Set(event.kline.quote_asset_volume);
        kline_model.open_ts = Set(event.kline.open_time);
        kline_model.close_ts = Set(event.kline.close_time);
        kline_model.trade_count = Set(event.kline.number_of_trades);
        kline_model.taker_buy_volume = Set(event.kline.taker_buy_base_asset_volume);
        kline_model.taker_buy_amount = Set(event.kline.taker_buy_quote_asset_volume);
        kline_model.average = Set((amount / volume).to_string());
        kline_model
    }
}

impl From<KlineSummary> for ActiveModel {
    fn from(summary: KlineSummary) -> Self {
        let amount = Decimal::from_str(&summary.quote_asset_volume)
            .inspect_err(|e| log::error!("get amount from kline summary: {:?}", e))
            .unwrap();
        let volume = Decimal::from_str(&summary.volume)
            .inspect_err(|e| log::error!("get volume from kline summary: {:?}", e))
            .unwrap();

        let mut kline = Self::new();
        kline.symbol = NotSet;
        kline.open = Set(summary.open);
        kline.close = Set(summary.close);
        kline.high = Set(summary.high);
        kline.low = Set(summary.low);
        kline.volume = Set(summary.volume.clone());
        kline.amount = Set(summary.quote_asset_volume.clone());
        kline.open_ts = Set(summary.open_time);
        kline.close_ts = Set(summary.close_time);
        kline.trade_count = Set(summary.number_of_trades);
        kline.taker_buy_volume = Set(summary.taker_buy_base_asset_volume);
        kline.taker_buy_amount = Set(summary.taker_buy_quote_asset_volume);
        kline.average = Set((amount / volume).to_string());
        kline
    }
}
