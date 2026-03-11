use std::env;
use std::net::TcpStream;
use std::io::{self, Write, Read};
use std::thread;

fn main() {

    let args: Vec<String> = env::args().collect();

    if args.len() != 3 {
        println!("Uso: cargo run <ip> <puerto>");
        return;
    }

    let host = &args[1];
    let puerto = &args[2];
    let direccion = format!("{}:{}", host, puerto);

    println!("\nIntentando conectar a {}...", direccion);

    //intentamos conectar el cliente con el servidor
    let mut hilo_cliente_servidor= match TcpStream::connect(direccion) {
	Ok(s) => {
            println!("Conectado al servidor :)");
            s
	}
	Err(e) => {
            println!("Error: {}", e);
            return;
	}
    };

    //como Rust no permite que dos hilos usen el mismo objeto al mismo tiempo, duplicamos
    //el descriptor (FD) del socket
    let hilo_servidor_cliente = hilo_cliente_servidor.try_clone().expect("No se pudo clonar el stream");

    //creamos un nuevo hilo para recibir mensajes del servidor
    iniciar_receptor(hilo_servidor_cliente);
    
    //mantenemos vivo al usuario para que escriba todo lo que quiera
    loop {

	let mut input = String::new();

	io::stdin()
            .read_line(&mut input)
            .expect("Error leyendo input");
	
	hilo_cliente_servidor
	    .write_all(input.as_bytes())
	    .unwrap();
	
    }
    
}


fn iniciar_receptor(mut hilo_servidor_cliente: TcpStream) {

    thread::spawn(move || {
        let mut buffer = [0; 512];
        loop {
            match hilo_servidor_cliente.read(&mut buffer) {
                Ok(n) if n > 0 => {
                    let mensaje = String::from_utf8_lossy(&buffer[..n]);
                    println!("{}", mensaje);
                }
		
                Ok(_) => {
                    println!("Servidor cerró la conexión");
                    break;
                }
		
                Err(e) => {
                    println!("Error recibiendo datos: {}", e);
                    break;
                }

            }
	    
        }

    });
    
}
