#[macro_use] extern crate rocket;
use std::collections::HashMap;

use rocket::State;

struct AppState {
    static_files: HashMap<&'static str, &'static str>
}

#[get("/")]
fn index(state: &State<AppState>) -> &'static str {
    state.static_files["index.html"]
}

#[launch]
fn rocket() -> _ {
    let mut static_files = HashMap::new();
    static_files.insert("index.html", "<p>Hello world!</p>");

    rocket::build().mount("/", routes![index]).manage(AppState {
        static_files
    })
}

