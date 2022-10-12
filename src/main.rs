#[macro_use]
extern crate rocket;
use rocket::response::content;
use rocket::serde::json::{Json};
use rocket::serde::{Deserialize};
use rocket::State;
use std::fs;
use std::vec;
use std::sync::Mutex;

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
    bookmarks: Mutex<Vec<Bookmark>>
}

// API routes
#[post("/bookmark", format = "json", data = "<bookmark>")]
async fn post_bookmark(state: &State<AppState>, bookmark: Json<JSONBookmark<'_>>) {
    let mut bookmarks = state.bookmarks.lock().unwrap();
       bookmarks.push(Bookmark {
       title: bookmark.title.to_string(),
       href: bookmark.href.to_string()
   });
}

// Static routes
#[get("/")]
async fn index_html(state: &State<AppState>) -> content::RawHtml<&'static str> {
    let mut result : String = String::new();
    let bookmarks = state.bookmarks.lock().unwrap();

    for item in bookmarks.iter() {
        result = format!("{}\n<a href=\"{}\">{}</a>", result, item.href, item.title);
    }

    let replaced_html = state.index_html.replace("<!--@bookmarks@-->", &format!("{}\n<!--@bookmarks@-->", &result));

    content::RawHtml(replaced_html).clone()
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
    rocket::build()
        .mount("/", routes![index_html, main_js, style_css, post_bookmark])
        .manage(AppState {
            index_html: fs::read_to_string("static/index.html").unwrap(),
            main_js: fs::read_to_string("static/main.js").unwrap(),
            style_css: fs::read_to_string("static/style.css").unwrap(),
            bookmarks: Mutex::new(Vec::new())
        })
}
