use std::sync::Arc;

use sea_orm::DbConn;
use tokio_tungstenite::tungstenite::Message;

use super::clients::Clients;
use super::Config;
use crate::cexs::binance::market::MarketClient;
use crate::service::prelude::*;

pub struct Server {
    pub(crate) cfg: Config,
    pub(crate) db: Arc<DbConn>,
    pub(crate) client: Arc<MarketClient>,
}

impl Server {
    pub fn new(cfg: Config, client: Clients) -> Self {
        Self {
            cfg,
            db: client.db,
            client: client.client,
        }
    }

    pub async fn init(&mut self) {
        let mut_client = Arc::get_mut(&mut self.client).unwrap();
        mut_client.set_handler(move |msg: Message| todo!(r#"handle message"#));
        // (, move |message: Message| {
        //     let db = conn.clone();
        //     async move {
        //         match message {
        //             Message::Text(text) => {
        //                 let model = kline::ActiveModel::from_event(text);
        //                 let table_name =
        //                     get_kline_table_name(&model.symbol.clone().unwrap(), model.interval);
        //                 let query = Query::insert().into_table(table_name);
        //                 Kline::insert(klines).await.unwrap();
        //             }
        //             _ => {}
        //         }
        //     };
        //     todo!()
        // })
        mut_client.set_reconnect_callback(move || todo!(r#"callback on reconnect"#));
        init_service(self.db.clone(), self.client.clone(), &self.cfg.service).await;
    }

    pub async fn run(&self) {}
}
