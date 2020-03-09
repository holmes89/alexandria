use serde::Serialize;
use diesel::types::Uuid;
use diesel::types::Timestamp;

#[derive(Queryable, Serialize, Insertable)]
pub struct JournalEntry {
    id: Uuid,
    content: String,
    created: Timestamp
}
