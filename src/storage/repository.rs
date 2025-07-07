use super::mysql::models::kline_1m::sea_orm::{ActiveModel, KlineModel};
use std::error::Error;
pub type Result<T> = std::result::Result<T, Box<dyn Error>>;

pub trait KlineRepository {
    async fn find_first(&self, symbol: &str) -> Result<Option<KlineModel>>;
    async fn find_last(&self, symbol: &str) -> Result<Option<KlineModel>>;
    async fn insert(&self, models: Vec<ActiveModel>) -> Result<Vec<KlineModel>>;
}
