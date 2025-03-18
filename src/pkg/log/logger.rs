use super::config::Config;
use flexi_logger::style;
use flexi_logger::{Age, Cleanup, Criterion, Duplicate, FileSpec, Logger, Naming, WriteMode};

pub fn init_logger(config: &Config) {
    // 初始化 flexi_logger
    let mut logger = Logger::try_with_str(&config.level)
        .unwrap()
        .log_to_file(
            FileSpec::default()
                .directory(&config.directory)
                .basename(&config.basename)
                .suffix(&config.suffix),
        )
        // 添加格式化设置，包含时间戳
        .format_for_files(|w, now, record| {
            write!(
                w,
                "{} [{}] {} - {}",
                now.format("%Y-%m-%d %H:%M:%S%.3f"),
                record.level(),
                record.target(),
                &record.args()
            )
        })
        // 为标准错误输出也添加格式
        .format_for_stderr(|w, now, record| {
            let level = record.level();
            write!(
                w,
                "{} [{}] {} - {}",
                now.format("%Y-%m-%d %H:%M:%S%.3f"),
                style(level).paint(level.to_string()),
                record.target(),
                &record.args()
            )
        });

    // 设置写入模式
    logger = match config.write_mode.as_str() {
        "buffer_and_flush" => logger.write_mode(WriteMode::BufferAndFlush),
        "direct" => logger.write_mode(WriteMode::Direct),
        _ => logger.write_mode(WriteMode::BufferAndFlush),
    };

    // 设置轮转策略
    logger = match config.rotate.criterion.as_str() {
        "age" => {
            let age = match config.rotate.age.as_str() {
                "day" => Age::Day,
                "hour" => Age::Hour,
                _ => Age::Day,
            };
            logger.rotate(
                Criterion::Age(age),
                match config.rotate.naming.as_str() {
                    "numbers" => Naming::Numbers,
                    "timestamps" => Naming::Timestamps,
                    _ => Naming::Numbers,
                },
                match config.rotate.cleanup.as_str() {
                    "keep_log_files" => Cleanup::KeepLogFiles(config.rotate.keep),
                    _ => Cleanup::KeepLogFiles(config.rotate.keep),
                },
            )
        }
        _ => logger,
    };

    // 设置日志输出到 stderr
    logger = match config.duplicate_to_stderr.as_str() {
        "all" => logger.duplicate_to_stderr(Duplicate::All),
        "trace" => logger.duplicate_to_stderr(Duplicate::Trace),
        "debug" => logger.duplicate_to_stderr(Duplicate::Debug),
        "info" => logger.duplicate_to_stderr(Duplicate::Info),
        "warn" => logger.duplicate_to_stderr(Duplicate::Warn),
        "error" => logger.duplicate_to_stderr(Duplicate::Error),
        _ => logger.duplicate_to_stderr(Duplicate::Info),
    };

    // 启动 logger
    logger
        .start()
        .unwrap_or_else(|e| panic!("Logger initialization failed with {}", e));

    // 使用不同级别的日志宏
    log::trace!("This is a trace message");
    log::debug!("This is a debug message");
    log::info!("This is an info message");
    log::warn!("This is a warning message");
    log::error!("This is an error message");
}
