use std::env;
use ::axum::{routing::get, Router};

#[tokio::main]
async fn main() {
    let args: Vec<String> = env::args().collect();

    // Default host and port values if not provided as command-line arguments
    let mut host = "127.0.0.1";
    let mut port = 1080;

    if args.len() > 1 {
        host = &args[1];
    }

    if args.len() > 2 {
        port = args[2].parse().expect("Invalid port number");
    }

    let endpoint = format!("{}:{}", host, port);

    let app: Router = Router::new()
        .route("/", get(handler))
        .route("/health", get(ok));

    let listener = tokio::net::TcpListener::bind(endpoint).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

async fn handler() -> &'static str {
    "Hello from Rust server"
}

async fn ok() -> &'static str {
    "OK"
}