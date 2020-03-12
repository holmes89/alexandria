use crate::schema::journal_entry;
use chrono::NaiveDateTime;
use serde::Serialize;
use uuid::Uuid;

#[derive(Queryable, Serialize, Insertable)]
#[table_name = "journal_entry"]
pub struct JournalEntry {
    id: Uuid,
    content: String,
    created: NaiveDateTime,
}
