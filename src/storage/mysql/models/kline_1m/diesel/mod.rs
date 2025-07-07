pub mod model;
pub mod schema;

pub use model::Model as KlineModel;
pub use schema::klines;
pub use schema::klines::dsl::*;
