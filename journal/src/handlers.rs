use crate::db;
use crate::models::JournalEntry;
use rocket_contrib::json::{Json, JsonValue};
use rocket_contrib::uuid::Uuid;

#[get("/entry")]
pub fn find_all_entries(conn: db::DbConn) -> JsonValue {
    let entries = db::find_all(conn);
    json!(entries)
}

#[get("/entry/<id>")]
pub fn find_entry_by_id(id: Uuid, conn: db::DbConn) -> JsonValue {
    let entry = db::find_by_id(id.into_inner(), conn);
    json!(entry)
}

#[post("/entry", format = "json", data = "<entry>")]
pub fn create_entry(entry: Json<JournalEntry>, conn: db::DbConn) -> JsonValue {
    let resp = db::create(&entry.0, conn);
    json!(resp)
}
