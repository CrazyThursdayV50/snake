use std::sync::Arc;

use super::config::Config;
use crate::pkg::args;
use crate::pkg::config::prelude::*;
use crate::pkg::log::prelude::*;
use crate::pkg::orm::prelude::*;
use crate::service::prelude::*;

pub async fn run() {
    let args = args::parse();
    let mut cfg = load_config::<Config>(&args.config);
    cfg.orm.update_by_env();

    init_logger(&cfg.log);

    let conn = connect(&cfg.orm).await.unwrap();
    let conn = Arc::new(conn);

    init_service(&conn, &cfg.service).await;
}
