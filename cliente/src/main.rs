mod conexion;
mod usuario;

fn main() {

    //obtenemos la direccion del servidor al que nos conectaremos
    let direccion = usuario::obtener_direccion();
    println!("Intentando conectar a {}", direccion);

    //intentamos conectar el cliente con el servidor
    let hilo_cliente_servidor= conexion::conectar(&direccion);

    //como Rust no permite que dos hilos usen el mismo objeto al mismo tiempo, duplicamos
    //el descriptor (FD) del socket para poder recibir del servidor
     conexion::leer_mensajes_servidor(hilo_cliente_servidor.try_clone().unwrap());
    
    //mantenemos vivo al usuario para que escriba todo lo que quiera
    usuario::escucha_usuario(hilo_cliente_servidor);
}
