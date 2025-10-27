pub mod models;

use crate::schema::*;
use diesel::{
    dsl::sql,
    expression::{SqlLiteral, TypedExpressionType},
    prelude::*,
};
use futures::stream::{self, Stream, StreamExt};
use lightspeed::models::OrderWrapper;
use models::*;

use chrono::Local;
use diesel_async::{
    AsyncPgConnection as Connection, RunQueryDsl,
    pooled_connection::{
        AsyncDieselConnectionManager as ConnectionManager, PoolError,
        bb8::{self, Pool},
    },
    scoped_futures::ScopedFutureExt,
};
use thiserror::Error;
use uuid::Uuid;

#[derive(Debug, Error)]
pub enum DatabaseError {
    #[error("unable to open database pool: '{0}'")]
    Pool(#[from] PoolError),
    #[error("bb8 error: '{0}'")]
    Bb8(#[from] bb8::RunError),
    #[error("database error: '{0}'")]
    Diesel(#[from] diesel::result::Error),
}

pub type ConnectionPool = Pool<Connection>;

pub type Result<T> = std::result::Result<T, DatabaseError>;

pub async fn establish_connection(url: &str) -> Result<Pool<Connection>> {
    Ok(Pool::builder()
        .build(ConnectionManager::<Connection>::new(url))
        .await?)
}

fn default<ST>() -> SqlLiteral<ST>
where
    ST: TypedExpressionType,
{
    sql("DEFAULT")
}

pub async fn create_order(pool: &Pool<Connection>, incoming: &OrderWrapper) -> Result<Order> {
    let mut conn = pool.get().await?;
    let data = serde_json::to_value(incoming.clone()).expect("valid incoming order data");

    Ok(diesel::insert_into(orders::table)
        .values((
            orders::incoming_order.eq(data),
            orders::lightspeed_order_id.eq(incoming.order.id as i32),
            orders::lightspeed_order_number.eq(incoming.order.number.to_string()),
        ))
        .returning(Order::as_returning())
        .get_result(&mut conn)
        .await?)
}

pub async fn link_dhl_draft(pool: &Pool<Connection>, draft_id: &Uuid, order_id: i32) -> Result<()> {
    let mut conn = pool.get().await?;

    diesel::update(orders::table)
        .filter(orders::lightspeed_order_id.eq(order_id))
        .set(orders::dhl_draft_id.eq(draft_id))
        .execute(&mut conn)
        .await?;

    Ok(())
}

pub async fn set_processed(pool: &Pool<Connection>, draft_id: &Uuid) -> Result<()> {
    let mut conn = pool.get().await?;
    diesel::update(orders::table)
        .filter(orders::dhl_draft_id.eq(draft_id))
        .set((
            orders::processed_at.eq(Local::now().naive_local()),
            orders::updated_at.eq(default()),
        ))
        .execute(&mut conn)
        .await?;
    Ok(())
}

pub async fn set_cancelled(pool: &Pool<Connection>, draft_id: &Uuid) -> Result<()> {
    let mut conn = pool.get().await?;
    diesel::update(orders::table)
        .filter(orders::dhl_draft_id.eq(draft_id))
        .set((
            orders::cancelled_at.eq(Local::now().naive_local()),
            orders::updated_at.eq(default()),
        ))
        .execute(&mut conn)
        .await?;
    Ok(())
}

pub async fn set_shipment_id(
    pool: &Pool<Connection>,
    draft_id: &Uuid,
    shipment_id: Option<Uuid>,
) -> Result<()> {
    let mut conn = pool.get().await?;
    diesel::update(orders::table)
        .filter(orders::dhl_draft_id.eq(draft_id))
        .set(orders::dhl_shipment_id.eq(shipment_id))
        .execute(&mut conn)
        .await?;

    Ok(())
}

pub async fn incr_poll_count(pool: &Pool<Connection>, order_id: i32) -> Result<()> {
    let mut conn = pool.get().await?;
    diesel::update(orders::table)
        .filter(orders::lightspeed_order_id.eq(order_id))
        .set((
            orders::poll_count.eq(orders::poll_count + 1),
            orders::updated_at.eq(default()),
        ))
        .execute(&mut conn)
        .await?;
    Ok(())
}

pub async fn unprocessed_count(pool: &Pool<Connection>) -> Result<i64> {
    let mut conn = pool.get().await?;
    Ok(orders::table
        .count()
        .filter(
            orders::processed_at
                .is_null()
                .and(orders::cancelled_at.is_null()),
        )
        .get_result(&mut conn)
        .await?)
}

pub async fn processed_count(pool: &Pool<Connection>) -> Result<i64> {
    let mut conn = pool.get().await?;
    Ok(orders::table
        .count()
        .filter(orders::processed_at.is_not_null())
        .get_result(&mut conn)
        .await?)
}

pub async fn get_unprocessed_stream(
    pool: &Pool<Connection>,
) -> Result<impl Stream<Item = QueryResult<Order>>> {
    let pool = pool.clone();

    let order_ids: Vec<i32> = {
        let mut conn = pool.get().await?;
        orders::table
            .filter(
                orders::processed_at
                    .is_null()
                    .and(orders::cancelled_at.is_null()),
            )
            .order(orders::created_at.asc())
            .select(orders::id)
            .load(&mut conn)
            .await?
    };

    Ok(stream::iter(order_ids).then(move |id| {
        let pool = pool.clone();
        async move {
            let Ok(mut conn) = pool.get().await else {
                // this is not ideal, i'm not sure if we can return a custom error here
                return Err(diesel::result::Error::BrokenTransactionManager);
            };
            orders::table
                .find(id)
                .select(Order::as_select())
                .first(&mut conn)
                .await
        }
        .scope_boxed()
    }))
}
