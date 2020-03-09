use rocket_contrib::databases::diesel;

use crate::models::{JournalEntry};
use crate::schema::{journalEntry};

use diesel;
use diesel::prelude::*;
use diesel::r2d2::{self, ConnectionManager};
use diesel::pg::PgConnection;

type Pool = r2d2::Pool<ConnectionManager<PgConnection>>;

pub struct DbConn(pub r2d2::PooledConnection<ConnectionManager<PgConnection>>);

pub fn find_by_id(id: &str, pool: Pool) -> Option<JournalEntry> {
    let conn: &PgConnection = &pool.get().unwrap();
    let entry = journalEntry::table
        .find(id)
        .first(conn)
        .map_err(|err| eprintln!("journalEntry::find_one: {}", err))
        .ok()?;
    Some(entry)
}

pub fn find_all(pool: Pool) -> Option<Vec<JournalEntry>> {
    let conn: &PgConnection = &pool.get().unwrap();
    journalEntry::table
        .load::<JournalEntry>(conn)
        .map_err(|err| eprintln!("journalEntry::find_all: {}", err))
        .ok()
}

pub fn create(entry: &JournalEntry, pool: Pool) -> Option<JournalEntry> {
    let conn: &PgConnection = &pool.get().unwrap();
    let entry: JournalEntry = diesel::insert(&comment)
       .into(comments::journalEntry)
       .get_result(&connection)
       .expect("Error saving new post");    
    Some(entry)
}
