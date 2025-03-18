use std::sync::Arc;

use sea_orm::DbConn;

use super::config::Config;
use super::on_load::prelude::*;

pub async fn init(db: &Arc<DbConn>, cfg: &Config) {
    init_kline_schema(db, cfg).await;
}
