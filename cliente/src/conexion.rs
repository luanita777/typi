use std::net::TcpStream;
use std::thread;
use std::io::Read;

//nos conecta al servidor
pub fn conectar(direccion: &str) -> TcpStream {
    TcpStream::connect(direccion)
        .expect("No se pudo conectar al servidor")
}


//lee los mensajes del servidor
pub fn leer_mensajes_servidor(mut hilo_servidor_cliente: TcpStream) {

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
