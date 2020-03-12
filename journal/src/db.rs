use crate::models::JournalEntry;
use crate::schema::journal_entry;

use ::diesel;
use diesel::pg::PgConnection;
use diesel::prelude::*;
use diesel::r2d2::{self, ConnectionManager};
use rocket::http::Status;
use rocket::request::{self, FromRequest};
use rocket::{Outcome, Request, State};
use std::env;
use std::ops::Deref;

pub fn init_pool() -> Pool {
    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    let manager = ConnectionManager::<PgConnection>::new(database_url);
    r2d2::Pool::builder()
        .build(manager)
        .expect("failed to create pool")
}

type Pool = r2d2::Pool<ConnectionManager<PgConnection>>;

pub struct DbConn(pub r2d2::PooledConnection<ConnectionManager<PgConnection>>);

impl<'a, 'r> FromRequest<'a, 'r> for DbConn {
    type Error = ();

    fn from_request(request: &'a Request<'r>) -> request::Outcome<DbConn, Self::Error> {
        let pool = request.guard::<State<Pool>>()?;
        match pool.get() {
            Ok(conn) => Outcome::Success(DbConn(conn)),
            Err(_) => Outcome::Failure((Status::ServiceUnavailable, ())),
        }
    }
}

impl Deref for DbConn {
    type Target = PgConnection;

    fn deref(&self) -> &Self::Target {
        &self.0
    }
}

pub fn find_by_id(id: uuid::Uuid, conn: DbConn) -> Option<JournalEntry> {
    let entry = journal_entry::table
        .find(id)
        .first(&*conn)
        .map_err(|err| eprintln!("journalEntry::find_one: {}", err))
        .ok()?;
    Some(entry)
}

pub fn find_all(conn: DbConn) -> Option<Vec<JournalEntry>> {
    journal_entry::table
        .load::<JournalEntry>(&*conn)
        .map_err(|err| eprintln!("journalEntry::find_all: {}", err))
        .ok()
}

pub fn create(entry: &JournalEntry, conn: DbConn) -> Option<JournalEntry> {
    let entry: JournalEntry = diesel::insert_into(journal_entry::table)
        .values(entry)
        .get_result(&*conn)
        .expect("Error saving new post");
    Some(entry)
}
