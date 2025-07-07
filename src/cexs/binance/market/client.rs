use crate::storage::mysql::models::kline_1m;
use crate::storage::mysql::models::kline_1m::*;

use super::super::config::Config as BinanceConfig;
use binance_spot_connector_rust::market::klines::Klines;
use binance_spot_connector_rust::{
    http::{request::Request, Credentials},
    hyper::BinanceHttpClient,
    market_stream::kline::KlineStream,
    tokio_tungstenite::{BinanceWebSocketClient, WebSocketState},
};

use futures_util::StreamExt;
use hyper::client::HttpConnector;
use hyper_tls::HttpsConnector;
use log;
use tokio::{
    net::TcpStream,
    time::{sleep, Duration},
};
use tokio_tungstenite::{
    tungstenite::{protocol::Message, Error},
    MaybeTlsStream,
};
// 存储当前订阅信息，用于断线重连时重新订阅
struct SubscriptionInfo {
    symbol: String,
    kline_type: kline_1m::KlineInterval,
}

pub struct Client {
    conn: BinanceHttpClient<HttpsConnector<HttpConnector>>,
    ws: WebSocketState<MaybeTlsStream<TcpStream>>,
    handler: Option<Box<dyn Fn(Message) + Send + Sync + 'static>>,
    // 存储配置信息，用于重连
    config: BinanceConfig,
    // 存储订阅信息
    subscriptions: Vec<SubscriptionInfo>,
    // 重连回调函数
    reconnect_callback: Option<Box<dyn Fn() + Send + 'static>>,
    // 重连间隔(毫秒)
    reconnect_interval: u64,
}

impl Client {
    pub async fn new(cfg: &BinanceConfig) -> Self {
        let credentials = {
            if cfg.is_hmac {
                Credentials::from_hmac(cfg.api_key.to_string(), cfg.secret_key.to_string())
            } else {
                Credentials::from_ed25519(cfg.api_key.to_string(), cfg.secret_key.to_string())
            }
        };

        let (ws, _) = BinanceWebSocketClient::connect_async_default()
            .await
            .unwrap();

        Self {
            conn: BinanceHttpClient::default().credentials(credentials),
            ws,
            handler: None,
            config: cfg.clone(),
            subscriptions: Vec::new(),
            reconnect_callback: None,
            reconnect_interval: 5000, // 默认5秒重连
        }
    }

    pub fn set_handler<F>(&mut self, handler: F)
    where
        F: Fn(Message) + Send + Sync + 'static,
    {
        self.handler = Some(Box::new(handler));
    }

    // 设置重连回调函数
    pub fn set_reconnect_callback<F>(&mut self, callback: F)
    where
        F: Fn() + Send + 'static,
    {
        self.reconnect_callback = Some(Box::new(callback));
    }

    // 设置重连间隔时间
    pub fn set_reconnect_interval(&mut self, interval_ms: u64) {
        self.reconnect_interval = interval_ms;
    }

    // 重连WebSocket
    async fn reconnect(&mut self) -> Result<(), Error> {
        log::info!("正在尝试重新连接到Binance WebSocket...");

        // 创建新的WebSocket连接
        match BinanceWebSocketClient::connect_async_default().await {
            Ok((new_ws, _)) => {
                // 更新WebSocket连接
                self.ws = new_ws;

                // 重新订阅之前的所有订阅
                for subscription in &self.subscriptions {
                    self.ws
                        .subscribe(vec![&KlineStream::new(
                            &subscription.symbol,
                            subscription.kline_type.into(),
                        )
                        .into()])
                        .await;
                }

                // 调用重连回调函数
                if let Some(callback) = &self.reconnect_callback {
                    callback();
                }

                log::info!("Binance WebSocket重连成功，已恢复所有订阅");
                Ok(())
            }
            Err(e) => {
                log::error!("Binance WebSocket重连失败: {:?}", e);
                Err(Error::ConnectionClosed)
            }
        }
    }

    pub async fn subscribe_kline(&mut self, symbol: &str, kline_type: KlineInterval) {
        // 存储订阅信息
        self.subscriptions.push(SubscriptionInfo {
            symbol: symbol.to_string(),
            kline_type,
        });

        self.ws
            .subscribe(vec![&KlineStream::new(symbol, kline_type.into()).into()])
            .await;
    }

    pub async fn handle_kline(&mut self) {
        let stream = self.ws.as_mut();
        if let Some(message_result) = stream.next().await {
            match message_result {
                Ok(message) => {
                    // 处理正常消息
                    match self.handler.as_ref() {
                        Some(handler) => handler(message),
                        None => {}
                    }
                }

                Err(e) => {
                    // 处理连接错误
                    log::error!("WebSocket连接错误: {:?}", e);
                    // 尝试重连
                    let mut reconnect_attempts = 0;
                    let max_attempts = 10; // 最大重试次数

                    while reconnect_attempts < max_attempts {
                        reconnect_attempts += 1;
                        log::info!("尝试重连 {}/{}", reconnect_attempts, max_attempts);

                        // 等待一段时间后尝试重连
                        sleep(Duration::from_millis(self.reconnect_interval)).await;

                        match self.reconnect().await {
                            Ok(_) => {
                                // 重连成功
                                break;
                            }
                            Err(_) => {
                                // 重连失败，继续尝试
                                log::warn!("重连尝试 {}/{} 失败", reconnect_attempts, max_attempts);
                                // 指数退避增加重连间隔
                                sleep(Duration::from_millis(
                                    self.reconnect_interval * reconnect_attempts as u64 / 2,
                                ))
                                .await;
                            }
                        }
                    }

                    if reconnect_attempts >= max_attempts {
                        log::error!("达到最大重连尝试次数，无法恢复连接");
                    }
                }
            }
        }
    }

    // 持续监听K线数据，自动处理断线重连
    pub async fn start_kline_stream(&mut self) {
        loop {
            self.handle_kline().await;
        }
    }

    pub async fn get_klines_from(
        &self,
        symbol: &str,
        start_time: i64,
        stop_time: Option<i64>,
        kline_type: KlineInterval,
        include_current: bool,
    ) -> Vec<ActiveModel> {
        let mut start_time = start_time;
        if !include_current {
            start_time = kline_1m::next_time(start_time, kline_type);
        }

        let params =
            kline_params_by_start_time(symbol.to_string(), kline_type, start_time, stop_time);
        let request: Request = params.into();

        let result = self
            .conn
            .send(request)
            .await
            .unwrap()
            .into_body_str()
            .await
            .unwrap();

        ActiveModelArray::from(&result).into()
    }

    pub async fn get_klines_to(
        &self,
        symbol: &str,
        end_time: i64,
        stop_time: Option<i64>,
        kline_type: KlineInterval,
        include_current: bool,
    ) -> Vec<ActiveModel> {
        let mut end_time = end_time;
        if !include_current {
            end_time = kline_1m::last_time(end_time, kline_type);
        }

        let params = kline_params_by_end_time(symbol.to_string(), kline_type, end_time, stop_time);
        let request: Request = params.into();

        let result = self
            .conn
            .send(request)
            .await
            .unwrap()
            .into_body_str()
            .await
            .unwrap();

        ActiveModelArray::from(&result).into()
    }
}

pub fn kline_params_by_start_time(
    symbol: String,
    interval: KlineInterval,
    start_time: i64,
    stop_time: Option<i64>,
) -> Klines {
    let mut kline = Klines::new(&symbol, interval.into());
    kline = kline.start_time(start_time as u64);
    let mut end_time = get_end_time_by_start_time(start_time, interval);
    if let Some(stop_time) = stop_time {
        let stop_time = last_time(stop_time, interval);
        if stop_time < end_time {
            end_time = stop_time
        }
    }
    kline = kline.end_time(end_time as u64);
    kline
}

pub fn kline_params_by_end_time(
    symbol: String,
    interval: KlineInterval,
    end_time: i64,
    stop_time: Option<i64>,
) -> Klines {
    let mut klines = Klines::new(&symbol, interval.into());
    klines = klines.end_time(end_time as u64);
    let mut start_time = get_start_time_by_end_time(end_time, interval);
    if let Some(stop_time) = stop_time {
        if stop_time > start_time {
            start_time = stop_time
        }
    }
    klines = klines.start_time(start_time as u64);
    klines
}
