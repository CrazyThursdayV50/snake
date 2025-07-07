pub struct Config {
    pub is_testnet: bool,
    pub api_key: String,
    pub secret_key: String,
}

impl Default for Config {
    fn default() -> Self {
        Config {
            is_testnet: true,
            api_key: String::new(),
            secret_key: String::new(),
        }
    }
}
