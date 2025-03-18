use crate::pkg::log::prelude::*;
use crate::pkg::orm::prelude::*;
use crate::service::prelude::*;
use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Config {
    pub log: LogConfig,
    pub service: ServiceConfig,
    pub orm: OrmConfig,
}
