mod cexs;
mod models;
mod pkg;
mod server;
mod service;
mod storage;

#[tokio::main]
async fn main() {
    server::run().await;
}
