use super::OrmConfig;
use sea_orm::{ConnectOptions, Database, DatabaseConnection, DbErr};
use std::str::FromStr;

/// Connect to the database using the provided configuration
pub async fn connect(config: &OrmConfig) -> Result<DatabaseConnection, DbErr> {
    let level = log::LevelFilter::from_str(&config.sqlx_logging_level)
        .map_or(log::LevelFilter::Info, |k| k);

    let mut opt = ConnectOptions::new(config.url.clone());
    opt.max_connections(config.max_connections)
        .min_connections(config.min_connections)
        .connect_timeout(config.connect_timeout)
        .idle_timeout(config.idle_timeout)
        .max_lifetime(config.max_lifetime)
        .sqlx_logging(config.sqlx_logging)
        .sqlx_logging_level(level);

    let db = Database::connect(opt).await?;
    Ok(db)
}
