use ::axum::{routing::get, Router};

#[tokio::main]
async fn main() {
    let app: Router = Router::new()
        .route("/", get(handler))
        .route("/health", get(ok));

    let listener = tokio::net::TcpListener::bind("127.0.0.1:1080").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

async fn handler() -> &'static str {
    "Hello from Rust server"
}

async fn ok() -> &'static str {
    "OK"
}