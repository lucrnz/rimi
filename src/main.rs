#[macro_use]
extern crate rocket;
use rocket::response::content;
use rocket::serde::json::{Json};
use rocket::serde::{Deserialize};
use rocket::State;
use std::fs;
use std::vec;

#[derive(Deserialize)]
#[serde(crate = "rocket::serde")]
struct JSONBookmark<'r> {
    title: &'r str,
    href: &'r str
}

struct Bookmark {
    title: String,
    href: String
}

struct AppState {
    index_html: String,
    main_js: String,
    style_css: String,
    bookmarks: Vec<Bookmark>
}


// API routes
#[post("/bookmark", format = "json", data = "<bookmark>")]
async fn post_bookmark(state: &State<AppState>, bookmark: Json<JSONBookmark<'_>>) {
   state.bookmarks.push(Bookmark {
       title: bookmark.title.to_string(),
       href: bookmark.href.to_string()
   });
}

// Static routes
#[get("/")]
async fn index_html(state: &State<AppState>) -> content::RawHtml<&str> {
    content::RawHtml(&state.index_html)
}

#[get("/main.js")]
async fn main_js(state: &State<AppState>) -> content::RawJavaScript<&str> {
    content::RawJavaScript(&state.main_js)
}

#[get("/style.css")]
async fn style_css(state: &State<AppState>) -> content::RawCss<&str> {
    content::RawCss(&state.style_css)
}

#[launch]
pub fn rocket() -> _ {
    let bookmarks : Vec<Bookmark> = Vec::new();
    rocket::build()
        .mount("/", routes![index_html, main_js, style_css])
        .manage(AppState {
            index_html: fs::read_to_string("static/index.html").unwrap(),
            main_js: fs::read_to_string("static/main.js").unwrap(),
            style_css: fs::read_to_string("static/style.css").unwrap(),
            bookmarks
        })
}
