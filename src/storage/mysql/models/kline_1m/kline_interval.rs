use binance_spot_connector_rust::market::klines::KlineInterval;
use chrono::Duration;
use serde::{Deserialize, Deserializer};
use std::fmt;

/// K线类型枚举
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Interval {
    Min1,   // 1分钟
    Min3,   // 3分钟
    Min5,   // 5分钟
    Min15,  // 15分钟
    Min30,  // 30分钟
    Hour1,  // 1小时
    Hour2,  // 2小时
    Hour4,  // 4小时
    Hour6,  // 6小时
    Hour8,  // 8小时
    Hour12, // 12小时
    Day1,   // 1天
    Day3,   // 3天
    Week1,  // 1周
    Month1, // 1个月
}

impl Interval {
    /// 获取K线类型的字符串表示
    pub fn as_str(&self) -> &'static str {
        match self {
            Interval::Min1 => "1m",
            Interval::Min3 => "3m",
            Interval::Min5 => "5m",
            Interval::Min15 => "15m",
            Interval::Min30 => "30m",
            Interval::Hour1 => "1h",
            Interval::Hour2 => "2h",
            Interval::Hour4 => "4h",
            Interval::Hour6 => "6h",
            Interval::Hour8 => "8h",
            Interval::Hour12 => "12h",
            Interval::Day1 => "1d",
            Interval::Day3 => "3d",
            Interval::Week1 => "1w",
            Interval::Month1 => "1M",
        }
    }

    pub fn duration(&self) -> Duration {
        match self {
            Interval::Min1 => Duration::minutes(1),
            Interval::Min3 => Duration::minutes(3),
            Interval::Min5 => Duration::minutes(5),
            Interval::Min15 => Duration::minutes(15),
            Interval::Min30 => Duration::minutes(30),
            Interval::Hour1 => Duration::hours(1),
            Interval::Hour2 => Duration::hours(2),
            Interval::Hour4 => Duration::hours(4),
            Interval::Hour6 => Duration::hours(6),
            Interval::Hour8 => Duration::hours(8),
            Interval::Hour12 => Duration::hours(12),
            Interval::Day1 => Duration::days(1),
            Interval::Day3 => Duration::days(3),
            Interval::Week1 => Duration::weeks(1),
            Interval::Month1 => Duration::days(30),
        }
    }
}

impl Default for Interval {
    fn default() -> Self {
        Interval::Min1
    }
}

impl<'de> Deserialize<'de> for Interval {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        let s = String::deserialize(deserializer)?;
        match s.as_str() {
            "1m" => Ok(Interval::Min1),
            "3m" => Ok(Interval::Min3),
            "5m" => Ok(Interval::Min5),
            "15m" => Ok(Interval::Min15),
            "30m" => Ok(Interval::Min30),
            "1h" => Ok(Interval::Hour1),
            "2h" => Ok(Interval::Hour2),
            "4h" => Ok(Interval::Hour4),
            "6h" => Ok(Interval::Hour6),
            "8h" => Ok(Interval::Hour8),
            "12h" => Ok(Interval::Hour12),
            "1d" => Ok(Interval::Day1),
            "3d" => Ok(Interval::Day3),
            "1w" => Ok(Interval::Week1),
            "1M" => Ok(Interval::Month1),
            _ => Err(serde::de::Error::custom("invalid kline type")),
        }
    }
}

impl Into<KlineInterval> for Interval {
    fn into(self) -> KlineInterval {
        match self {
            Interval::Min1 => KlineInterval::Minutes1,
            Interval::Min3 => KlineInterval::Minutes3,
            Interval::Min5 => KlineInterval::Minutes5,
            Interval::Min15 => KlineInterval::Minutes15,
            Interval::Min30 => KlineInterval::Minutes30,
            Interval::Hour1 => KlineInterval::Hours1,
            Interval::Hour2 => KlineInterval::Hours2,
            Interval::Hour4 => KlineInterval::Hours4,
            Interval::Hour6 => KlineInterval::Hours6,
            Interval::Hour8 => KlineInterval::Hours8,
            Interval::Hour12 => KlineInterval::Hours12,
            Interval::Day1 => KlineInterval::Days1,
            Interval::Day3 => KlineInterval::Days3,
            Interval::Week1 => KlineInterval::Weeks1,
            Interval::Month1 => KlineInterval::Months1,
        }
    }
}

impl fmt::Display for Interval {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}", self.as_str())
    }
}
