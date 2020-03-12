#![feature(proc_macro_hygiene, decl_macro)]

#[macro_use]
extern crate rocket;

#[macro_use]
extern crate rocket_contrib;

#[macro_use]
extern crate diesel;

mod db;
mod handlers;
mod models;
mod schema;

use rocket_contrib::json::JsonValue;

#[catch(404)]
fn not_found() -> JsonValue {
    json!({
        "status": "error",
        "reason": "Resource was not found."
    })
}

pub fn rocket() -> rocket::Rocket {
    rocket::ignite()
        .mount(
            "/api",
            routes![handlers::find_all_entries, handlers::find_entry_by_id,],
        )
        .manage(db::init_pool())
        .register(catchers![not_found])
}
