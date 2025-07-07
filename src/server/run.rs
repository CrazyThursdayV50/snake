use super::clients::Clients;
use super::Config;
use super::Server;
use crate::pkg::args;
use crate::pkg::config::prelude::*;

pub async fn run() {
    // init
    let args = args::parse();
    let mut cfg = load_config::<Config>(&args.config);
    cfg.orm.update_by_env();

    let client = Clients::new(&cfg).await;
    let server = Server::new(cfg, client);
    todo!()
}
