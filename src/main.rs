#[macro_use] extern crate rocket;
use std::fs;
use rocket::State;
use rocket::response::content;

struct AppState {
    index_html: String,
    main_js: String,
    style_css: String
}

#[get("/")]
fn index_html(state: &State<AppState>) -> content::RawHtml<& str> {
    content::RawHtml(&state.index_html)
}

#[get("/main.js")]
fn main_js(state: &State<AppState>) -> content::RawJavaScript<& str> {
    content::RawJavaScript(&state.main_js)
}

#[get("/style.css")]
fn style_css(state: &State<AppState>) -> content::RawCss<& str> {
    content::RawCss(&state.style_css)
}

#[launch]
fn rocket() -> _ {
    rocket::build().mount("/", routes![index_html, main_js, style_css]).manage(AppState {
    index_html: fs::read_to_string("static/index.html").unwrap(),
    main_js: fs::read_to_string("static/main.js").unwrap(),
    style_css: fs::read_to_string("static/style.css").unwrap()
    })
}

