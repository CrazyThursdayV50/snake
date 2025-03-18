use serde::{Deserialize, Serialize};
use std::time::Duration;

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Config {
    pub url: String,
    pub max_connections: u32,
    pub min_connections: u32,
    pub connect_timeout: Duration,
    pub idle_timeout: Duration,
    pub max_lifetime: Duration,
    pub sqlx_logging: bool,
    pub sqlx_logging_level: String,
}

impl Default for Config {
    fn default() -> Self {
        Self {
            url: "mysql://{user}:{password}@{ip}:{port}/{database}".to_string(),
            max_connections: 10,
            min_connections: 2,
            connect_timeout: Duration::from_secs(10),
            idle_timeout: Duration::from_secs(600),
            max_lifetime: Duration::from_secs(1800),
            sqlx_logging: false,
            sqlx_logging_level: "info".to_string(),
        }
    }
}

use std::env;
const MYSQL_USER: &str = "MYSQL_USER";
const MYSQL_PASSWORD: &str = "MYSQL_PASSWORD";
const MYSQL_IP: &str = "MYSQL_IP";
const MYSQL_PORT: &str = "MYSQL_PORT";
const MYSQL_DEFAULT_DB: &str = "MYSQL_DEFAULT_DB";

impl Config {
    pub fn update_by_env(&mut self) {
        let user = env::var(MYSQL_USER).map_or("root".to_string(), |u| u);
        let password = env::var(MYSQL_PASSWORD).map_or("".to_string(), |u| u);
        let ip = env::var(MYSQL_IP).map_or("127.0.0.1".to_string(), |u| u);
        let port = env::var(MYSQL_PORT).map_or("3306".to_string(), |u| u);
        let db = env::var(MYSQL_DEFAULT_DB).map_or("test".to_string(), |u| u);

        self.url = self
            .url
            .replace("{user}", &user)
            .replace("{password}", &password)
            .replace("{ip}", &ip)
            .replace("{port}", &port)
            .replace("{database}", &db)
    }
}
