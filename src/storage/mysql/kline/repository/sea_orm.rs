use crate::storage::mysql::models::kline_1m::sea_orm::{ActiveModel, Kline, KlineModel};
use crate::storage::repository::{KlineRepository, Result};
use sea_orm::*;
use sea_orm::{DbConn, EntityTrait, QueryOrder};
use std::sync::Arc;

pub(super) struct Repository {
    db: Arc<DbConn>,
}

impl KlineRepository for Repository {
    async fn insert(&self, models: Vec<ActiveModel>) -> Result<Vec<KlineModel>> {
        Ok(Kline::insert_many(models)
            .exec_with_returning_many(self.db.as_ref())
            .await?)
    }

    async fn find_last(&self, symbol: &str) -> Result<Option<KlineModel>> {
        use crate::storage::mysql::models::kline_1m::sea_orm::Column::*;
        Ok(Kline::find()
            .filter(Symbol.eq(symbol))
            .order_by_desc(OpenTs)
            .one(self.db.as_ref())
            .await
            .map_err(|e| e.to_string())?)
    }

    async fn find_first(&self, symbol: &str) -> Result<Option<KlineModel>> {
        use crate::storage::mysql::models::kline_1m::sea_orm::Column::*;
        Ok(Kline::find()
            .filter(Symbol.eq(symbol))
            .order_by_asc(OpenTs)
            .one(self.db.as_ref())
            .await
            .map_err(|e| e.to_string())?)
    }
}
