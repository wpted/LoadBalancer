use std::io::{Read, Write};
use std::net::{TcpListener, TcpStream};

fn main() {
    println!("Hello, world!");
    const HOST: &str = "127.0.0.1";
    const PORT: &str = "1080";

    let end_point: String = HOST.to_owned() + ":" + PORT;

    let listener =  TcpListener::bind(end_point).unwrap();
    println!("Web Server is listening at port {}", PORT);

    for stream in listener.incoming() {
        let _stream = stream.unwrap();
        handle_connection(_stream);
    }
}

fn handle_connection(mut stream: TcpStream) {
    let mut buffer = [0; 1024];
    stream.read(&mut buffer).unwrap();
    println!("Request: {}", String::from_utf8_lossy(&buffer[..]));

    let response_contents = "Hello from Rust server!";
    let response = format!(
        "HTTP/1.1 200 OK\r\nContent-Length: {}\r\n\r\n{}",
        response_contents.len(),
        response_contents
    );
    stream.write(response.as_bytes()).unwrap();
    stream.flush().unwrap();
}
