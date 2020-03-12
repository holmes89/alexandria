use crate::schema::journal_entry;
use chrono::NaiveDateTime;
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Queryable, Serialize, Deserialize, Insertable)]
#[table_name = "journal_entry"]
pub struct JournalEntry {
    id: Option<Uuid>,
    content: String,
    created: Option<NaiveDateTime>,
}
