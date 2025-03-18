use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Config {
    pub level: String,
    pub directory: String,
    pub basename: String,
    pub suffix: String,
    pub write_mode: String,
    pub rotate: Rotate,
    pub duplicate_to_stderr: String,
}

#[derive(Debug, Deserialize)]
pub struct Rotate {
    pub criterion: String,
    pub age: String,
    pub naming: String,
    pub cleanup: String,
    pub keep: usize,
}

impl Config {
    pub fn default() -> Self {
        Config {
            level: "debug".to_string(),
            directory: "/tmp".to_string(),
            basename: "main".to_string(),
            suffix: "log".to_string(),
            write_mode: "buffer_and_flush".to_string(),
            rotate: Rotate::default(),
            duplicate_to_stderr: "all".to_string(),
        }
    }
}

impl Rotate {
    fn default() -> Self {
        Rotate {
            criterion: "age".to_string(),
            age: "day".to_string(),
            naming: "timestamp".to_string(),
            cleanup: "".to_string(),
            keep: 7,
        }
    }
}
