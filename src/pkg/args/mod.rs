use clap::Parser;

/// 命令行参数结构体
#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
pub struct Args {
    /// 配置文件路径
    #[arg(short, long, default_value = "config.yaml")]
    pub config: String,
}

pub fn parse() -> Args {
    Args::parse()
}
