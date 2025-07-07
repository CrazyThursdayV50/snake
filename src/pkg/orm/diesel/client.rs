use super::OrmConfig;
// use diesel::{
//     prelude::*,
//     r2d2::{ConnectionManager, Pool},
// };

// use diesel_async::pooled_connection::deadpool::Pool;
// use diesel_async::pooled_connection::AsyncDieselConnectionManager;
// use diesel_async::RunQueryDsl;
// use diesel_async::{AsyncMysqlConnection, RunQueryDsl};
// use std::env;

// fn connect(cfg: &OrmConfig) {
//     let config = AsyncDieselConnectionManager::<AsyncMysqlConnection>::new(cfg.url);

//     // let pool = Pool::builder(config).build()?;
//     let pool = Pool::builder(config)
//         // .max_lifetime(Some(cfg.max_lifetime))
//         // .idle_timeout(Some(cfg.idle_timeout))
//         // .connection_timeout(cfg.connect_timeout)
//         // .max_size(cfg.max_connections)
//         // .min_idle(Some(cfg.min_connections))
//         // .build(manager)
//         .build()
//         .unwrap();

//     pool
// }

use diesel::prelude::*;
use diesel::r2d2::{self, ConnectionManager};

pub type Pool = r2d2::Pool<ConnectionManager<MysqlConnection>>;

fn connect(cfg: &OrmConfig) -> Pool {
    let manager = ConnectionManager::<MysqlConnection>::new(&cfg.url);
    Pool::builder()
        .max_lifetime(Some(cfg.max_lifetime))
        .idle_timeout(Some(cfg.idle_timeout))
        .connection_timeout(cfg.connect_timeout)
        .max_size(cfg.max_connections)
        .min_idle(Some(cfg.min_connections))
        .build(manager)
        .unwrap()
}
