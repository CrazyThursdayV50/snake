use std::sync::Arc;

use log;
use sea_orm::{ConnectionTrait, DbConn, DbErr, Statement};
use sea_query::{ColumnDef, IndexCreateStatement, MysqlQueryBuilder, Table};

/// 判断表是否存在
async fn table_exists(db: &DbConn, table_name: &str) -> Result<bool, DbErr> {
    let stmt = Statement::from_string(
        db.get_database_backend(),
        format!(
            "SELECT COUNT(*) as count FROM information_schema.tables WHERE table_name = '{}'",
            table_name
        ),
    );

    let result = db.query_one(stmt).await?;
    let count: i64 = result.expect("查询失败").try_get("", "count")?;

    Ok(count > 0)
}

/// 创建K线表
///
/// 这是API模块唯一对外公开的函数，用于创建K线表
pub async fn create_kline_table(db: Arc<DbConn>, table_name: &str) -> Result<(), DbErr> {
    // 先检查表是否已经存在
    if table_exists(db.as_ref(), &table_name).await? {
        log::info!("表 {} 已存在，跳过创建", table_name);
        return Ok(());
    }

    // 创建表
    let mut table = Table::create();
    table
        .table(sea_query::Alias::new(table_name))
        .if_not_exists();

    // 添加列
    table.col(
        ColumnDef::new(sea_query::Alias::new("symbol"))
            .string_len(16)
            .not_null(),
    );
    table.col(
        ColumnDef::new(sea_query::Alias::new("open_ts"))
            .big_unsigned()
            .not_null()
            .primary_key(),
    );
    table.col(
        ColumnDef::new(sea_query::Alias::new("close_ts"))
            .big_unsigned()
            .not_null(),
    );
    table.col(
        ColumnDef::new(sea_query::Alias::new("open"))
            .string_len(40)
            .not_null(),
    );
    table.col(
        ColumnDef::new(sea_query::Alias::new("close"))
            .string_len(40)
            .not_null(),
    );
    table.col(
        ColumnDef::new(sea_query::Alias::new("low"))
            .string_len(40)
            .not_null(),
    );
    table.col(
        ColumnDef::new(sea_query::Alias::new("high"))
            .string_len(40)
            .not_null(),
    );
    table.col(
        ColumnDef::new(sea_query::Alias::new("average"))
            .string_len(40)
            .not_null(),
    );
    table.col(
        ColumnDef::new(sea_query::Alias::new("volume"))
            .string_len(40)
            .not_null(),
    );
    table.col(
        ColumnDef::new(sea_query::Alias::new("amount"))
            .string_len(40)
            .not_null(),
    );
    table.col(
        ColumnDef::new(sea_query::Alias::new("trade_count"))
            .big_unsigned()
            .null(),
    );
    table.col(
        ColumnDef::new(sea_query::Alias::new("taker_buy_volume"))
            .string_len(40)
            .null(),
    );
    table.col(
        ColumnDef::new(sea_query::Alias::new("taker_buy_amount"))
            .string_len(40)
            .null(),
    );
    table.col(
        ColumnDef::new(sea_query::Alias::new("created_at"))
            .timestamp()
            .not_null()
            .extra("DEFAULT CURRENT_TIMESTAMP".to_string()),
    );
    table.col(
        ColumnDef::new(sea_query::Alias::new("updated_at"))
            .timestamp()
            .not_null()
            .extra("DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP".to_string()),
    );

    // 添加索引
    let mut idx = IndexCreateStatement::new()
        .name(&format!("uk_{}_close_ts", table_name))
        .table(sea_query::Alias::new(table_name))
        .col(sea_query::Alias::new("close_ts"))
        .unique()
        .to_owned();
    table.index(&mut idx);

    // 执行创建表的语句
    let sql = table.build(MysqlQueryBuilder);
    db.execute(Statement::from_string(db.get_database_backend(), sql))
        .await
        .unwrap();

    log::info!("成功创建K线表: {}", table_name);
    Ok(())
}
