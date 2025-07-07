#[macro_export]
macro_rules! create_kline_model {
    ($mod_name:ident, $schema_name:literal, $schema_ident:ident) => {
        pub mod $mod_name {
            use chrono::NaiveDateTime;
            use rust_decimal::Decimal;
            use sea_orm::*;
            use serde::{Deserialize, Serialize};
            #[derive(
                Default, Clone, Debug, PartialEq, Eq, DeriveEntityModel, Serialize, Deserialize,
            )]
            #[sea_orm(table_name = $schema_name	)]
            pub struct Model {
                #[sea_orm(primary_key)]
                pub open_ts: i64,
                pub close_ts: i64,
                pub open: Decimal,
                pub close: Decimal,
                pub low: Decimal,
                pub high: Decimal,
                pub average: Decimal,
                pub volume: Decimal,
                pub amount: Decimal,
                pub created_at: NaiveDateTime,
                #[sea_orm(on_update = "CURRENT_TIMESTAMP")]
                pub updated_at: NaiveDateTime,
            }

            #[derive(DeriveIden)]
            pub enum $schema_ident {
                Table,
            }

            #[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
            pub enum Relation {}

            impl ActiveModelBehavior for ActiveModel {}
        }
    };
}
