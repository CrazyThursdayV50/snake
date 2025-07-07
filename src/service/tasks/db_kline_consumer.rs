use std::sync::Arc;

use sea_orm::{DbConn, EntityTrait};
use tokio::sync::mpsc;

use crate::storage::mysql::models::kline_1m::{ActiveModel, Kline};

pub struct Consumer {
    db: Arc<DbConn>,
    kline_tx: mpsc::Sender<Vec<ActiveModel>>,
    kline_rx: mpsc::Receiver<Vec<ActiveModel>>,
    done: mpsc::Receiver<()>,
    close: mpsc::Sender<()>,
}

impl Consumer {
    pub fn new(db: Arc<DbConn>) -> Self {
        let (kline_tx, kline_rx) = mpsc::channel::<Vec<ActiveModel>>(100);
        let (done_tx, done_rx) = mpsc::channel::<()>(0);
        Self {
            db,
            kline_tx,
            kline_rx,
            done: done_rx,
            close: done_tx,
        }
    }

    pub async fn run(&mut self) {
        loop {
            tokio::select! {
                Some(_) = self.done.recv() => {
                  log::warn!("db kline consumer exit");
                  break;
                },

                Some(klines) = self.kline_rx.recv() => {
                  if klines.len() > 0 {
                    Kline::insert_many(klines).exec(self.db.as_ref()).await.unwrap();
                  };
                },
            };
        }
    }

    pub async fn stop(&self) {
        self.close.send(()).await.unwrap();
    }

    pub async fn send_klines(&self, klines: Vec<ActiveModel>) {
        self.kline_tx.send(klines).await.unwrap();
    }
}
