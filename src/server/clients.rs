use std::sync::Arc;

use super::Config;
use crate::cexs::binance::market::MarketClient;
use crate::pkg::orm::prelude::*;
use sea_orm::DbConn;

pub(super) struct Clients {
    pub(crate) db: Arc<DbConn>,
    pub(crate) client: Arc<MarketClient>,
}

impl Clients {
    pub(crate) async fn new(cfg: &Config) -> Self {
        let db = connect(&cfg.orm).await.unwrap();
        let client = MarketClient::new(&cfg.service.cexs.binance).await;
        Clients {
            db: Arc::new(db),
            client: Arc::new(client),
        }
    }
}
