table! {
    journal_entry (id) {
        id -> Nullable<Uuid>,
        content -> Text,
        created -> Nullable<Timestamp>,
    }
}
