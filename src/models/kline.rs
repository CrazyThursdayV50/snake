use sea_orm::entity::prelude::*;
use serde::{Deserialize, Deserializer, Serialize};
use std::fmt;

/// K线类型枚举
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum KlineType {
    Min1,   // 1分钟
    Min5,   // 5分钟
    Min15,  // 15分钟
    Min30,  // 30分钟
    Hour1,  // 1小时
    Hour4,  // 4小时
    Day1,   // 1天
    Week1,  // 1周
    Month1, // 1个月
}

impl KlineType {
    /// 获取K线类型的字符串表示
    pub fn as_str(&self) -> &'static str {
        match self {
            KlineType::Min1 => "1m",
            KlineType::Min5 => "5m",
            KlineType::Min15 => "15m",
            KlineType::Min30 => "30m",
            KlineType::Hour1 => "1h",
            KlineType::Hour4 => "4h",
            KlineType::Day1 => "1d",
            KlineType::Week1 => "1w",
            KlineType::Month1 => "1M",
        }
    }
}

impl<'de> Deserialize<'de> for KlineType {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        let s = String::deserialize(deserializer)?;
        match s.as_str() {
            "1m" => Ok(KlineType::Min1),
            "5m" => Ok(KlineType::Min5),
            "15m" => Ok(KlineType::Min15),
            "30m" => Ok(KlineType::Min30),
            "1h" => Ok(KlineType::Hour1),
            "4h" => Ok(KlineType::Hour4),
            "1d" => Ok(KlineType::Day1),
            "1w" => Ok(KlineType::Week1),
            "1M" => Ok(KlineType::Month1),
            _ => Err(serde::de::Error::custom("invalid kline type")),
        }
    }
}

impl fmt::Display for KlineType {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}", self.as_str())
    }
}

/// 根据sql/kline.sql实现的K线结构体
#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Serialize, Deserialize)]
#[sea_orm(table_name = "kline")]
pub struct Model {
    #[sea_orm(primary_key)]
    pub open_ts: i64,
    pub close_ts: i64,
    pub open: String,
    pub close: String,
    pub low: String,
    pub high: String,
    pub average: String,
    pub volume: String,
    pub amount: String,
    #[sea_orm(auto_increment = false)]
    pub created_at: DateTimeWithTimeZone,
    #[sea_orm(auto_increment = false)]
    pub updated_at: DateTimeWithTimeZone,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {}

impl ActiveModelBehavior for ActiveModel {}

/// 获取表名
pub fn get_kline_table_name(symbol: &str, kline_type: KlineType) -> String {
    format!("kline_{}_{}", symbol.to_lowercase(), kline_type.as_str())
}
