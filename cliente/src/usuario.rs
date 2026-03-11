use std::env;
use std::io::{self, Write};
use std::net::TcpStream;

//obtenemos la direccion del servidor al que nos conectaremos
pub fn obtener_direccion() -> String {
    let args: Vec<String> = env::args().collect();
    if args.len() != 3 {
        println!("Uso: cargo run <ip> <puerto>");
        std::process::exit(1);
    }
    let host = &args[1];
    let puerto = &args[2];
    format!("{}:{}", host, puerto)
}

//se encarga de leer/escuchar todo lo que pone el usuario
pub fn escucha_usuario(mut stream: TcpStream) {
    loop {
        let mut input = String::new();
        io::stdin()
            .read_line(&mut input)
            .expect("Error leyendo input");
        stream.write_all(input.as_bytes()).unwrap();
    }
}
