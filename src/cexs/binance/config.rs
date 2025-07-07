use serde::Deserialize;

#[derive(Debug, Deserialize, Clone)]
pub struct Config {
    pub is_testnet: bool,
    pub api_key: String,
    pub secret_key: String,
    pub is_hmac: bool,
}

impl Default for Config {
    fn default() -> Self {
        Config {
            is_testnet: true,
            is_hmac: true,
            api_key: String::new(),
            secret_key: String::new(),
        }
    }
}
