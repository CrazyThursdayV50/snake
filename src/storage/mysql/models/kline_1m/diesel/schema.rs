diesel::table! {
    klines (open_ts) {
        #[max_length = 16]
        symbol -> Varchar,
        open_ts -> BigInt,
        close_ts -> BigInt,
        #[max_length = 40]
        open-> Varchar,
        #[max_length = 40]
        close-> Varchar,
        #[max_length = 40]
        low-> Varchar,
        #[max_length = 40]
        high-> Varchar,
        #[max_length = 40]
        average-> Varchar,
        #[max_length = 40]
        volume-> Varchar,
        #[max_length = 40]
        amount-> Varchar,
        trade_count-> Integer,
        #[max_length = 40]
        taker_buy_volume-> Varchar,
        #[max_length = 40]
        taker_buy_amount-> Varchar,
        created_at-> Datetime,
        updated_at-> Datetime,
    }
}
