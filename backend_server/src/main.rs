use std::io;
use std::io::{Read, Write};
use std::net::{TcpListener, TcpStream};

fn main() {
    println!("Hello, world!");
    const HOST: &str = "127.0.0.1";
    let mut port: String = String::new();
    // Use unwrap() to panic when the Result() is Err(error).
    // Use '?' to return from the call stack, enclosing the function with the error value.
    io::stdin().read_line(&mut port).unwrap();

    let endpoint: String = { HOST.to_owned() + ":" + &port };


    // Creates a TCP listener bound to the given endpoint. Panic if bind fails.
    let listener: TcpListener = TcpListener::bind(endpoint).unwrap();
    println!("Web Server is listening at port {}", port);

    // Connecting to any incoming connections
    for stream_result in listener.incoming() {
        match stream_result {
            Ok(stream) => {
                handle_connection(stream)
            }
            Err(err) => {
                println!("Error accepting connection: {}", err)
            }
        }
    }
}

fn handle_connection(mut stream: TcpStream) {
    let mut buffer: [u8; 1024] = [0; 1024];

    if let Err(err) = stream.read(&mut buffer) {
        println!("Error reading request to buffer: {}", err);
    }

    println!("Request: {}", String::from_utf8_lossy(&buffer[..]));

    let response_contents = "Hello from Rust server!";
    let response = format!(
        "HTTP/1.1 200 OK\r\nContent-Length: {}\r\n\r\n{}",
        response_contents.len(),
        response_contents
    );

    if let Err(err) = stream.write(response.as_bytes()) {
        println!("Error writing response: {}", err);
    }

    // Unwraps Result<()> to (), doesn't panic since there wouldn't be errors.
    stream.flush().unwrap();
}
