use serde::Deserialize;
use serde_yaml_ok;
use std::fs;

pub fn load_config<Config>(config_path: &str) -> Config
where
    Config: for<'de> Deserialize<'de> + Sized,
{
    let config_content = fs::read_to_string(config_path).expect("Failed to read config file");
    serde_yaml_ok::from_str(&config_content).expect("Failed to parse config file")
}
