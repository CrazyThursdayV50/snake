# Sea-ORM 和 InfluxDB 示例项目

这个项目展示了如何使用 Sea-ORM 和 InfluxDB 在 Rust 中进行数据库操作。

## 功能特性

- 使用 Sea-ORM 创建自定义表名的实体
- 动态创建 K 线表，支持不同的时间周期和交易对
- 从 SQL 文件创建表的迁移器
- 完全自定义表结构的表管理器

## 依赖

- Rust 1.70+
- MySQL 8.0+
- InfluxDB 2.0+

## 安装

```bash
# 克隆仓库
git clone https://github.com/yourusername/sea-orm-influxdb-example.git
cd sea-orm-influxdb-example

# 编译项目
cargo build
```

## 使用方法

### 创建 K 线表

```bash
# 使用默认数据库 URL 和交易对
cargo run --bin create_kline_tables

# 指定数据库 URL 和交易对
cargo run --bin create_kline_tables "mysql://username:password@localhost/dbname" "BTCUSDT"
```

这将创建以下表：
- kline_btcusdt_1m
- kline_btcusdt_5m
- kline_btcusdt_15m
- kline_btcusdt_1h
- kline_btcusdt_1d

### 从 SQL 文件创建表

```bash
# 使用默认数据库 URL 和 SQL 目录
cargo run --bin sql_migrator

# 指定数据库 URL 和 SQL 目录
cargo run --bin sql_migrator "mysql://username:password@localhost/dbname" "./sql"
```

### 测试 ORM 功能

```bash
# 显示帮助信息
cargo run --bin test_orm

# 创建 K 线表
cargo run --bin test_orm kline "mysql://username:password@localhost/dbname" "ETHUSDT"

# 从 SQL 文件创建表
cargo run --bin test_orm sql "mysql://username:password@localhost/dbname" "./sql"

# 创建自定义表
cargo run --bin test_orm custom "mysql://username:password@localhost/dbname" "my_custom_table"
```

## 运行测试

```bash
# 设置测试数据库 URL
export TEST_DB_URL="mysql://username:password@localhost/test_db"

# 运行测试
cargo test
```

## 代码示例

### 创建 K 线配置

```rust
// 创建 K 线配置
let cfg = KlineConfig::new(KlineType::Min1, "BTCUSDT");
// 表名将是 "kline_btcusdt_1m"
println!("表名: {}", cfg.get_table_name());

// 创建表
cfg.create_table(&db).await?;
```

### 使用表管理器创建自定义表

```rust
// 创建表管理器
let table_manager = TableManager::new();

// 设置自定义表名
table_manager.set_table_name::<Entity>("my_custom_kline");

// 创建表
table_manager.create_table::<Entity>(&db).await?;
```

### 从 SQL 文件创建表

```rust
// 创建 SQL 迁移器
let migrator = SqlMigrator::new(Path::new("./sql"));

// 从 SQL 文件创建表
migrator.create_tables_from_dir(&db).await?;
```

## 许可证

MIT 