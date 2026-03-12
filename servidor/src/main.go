package main

/*
   Creamos un servidor que acepta multiples clientes, tal que:
      1. Se lee el puerto desde la línea de comandos.
      2. Se abre un socket TCP en ese puerto.
      3. Esperamos conexiones de clientes.
      4. Cuando un cliente se conecta, se crea una goroutine
         para manejar esa conexión.
*/

import (
	"fmt"
	"net"
	"os"
)

func main() {

	//Si el usuario no pasa ningun puerto, terminamos
	if len(os.Args) != 2 {
		fmt.Println("Uso correcto del programa:")
		fmt.Println("go run servidor.go <puerto>")
		return
	}

	//Si el usuario sí paso un puerto lo guardamos
	var puerto string = os.Args[1]

	//Creamos una string de la direccion de nuestro servidor con nuestra ip
	//por omision y el puerto dado, esto solo hace :puerto
	var direccion string = ":" + puerto

	//Creamos el servidor, como Listen puede devolver dos valores
	//por eso se asigna a esos dos, además Listen se encarga de
	//darle sentido al texto que tiene dirección, pues interpreta
	//que el servidor va a estar en el puerto dado y como no le
	//pasamos una ip va a usar todas las que tenga disponibles
	//la computadora
	var listener net.Listener
	var err error

	listener, err = net.Listen("tcp", direccion)

	//si al crear el servidor sí hubo error, entonces err no es
	//nulo y terminamos el programa
	if err != nil {
		fmt.Println("Error al crear el servidor:")
		fmt.Println(err)
		return
	}

	//si err sí es nulo entonces el servidor se creó exitosamente
	fmt.Println("Servidor escuchando en el puerto", puerto)

	//Empezamos a aceptar clientes
	for {
		//Accept espera hasta que un cliente se conecte,
		//cuando alguien se conecta regresa la conexión
		//con ese cliente, si la conexion no es exitosa
		//enotnces err es distinto de null, pero no
		//terminamos el programam simplemente seguimos
		//aceptando conexiones

		var conexion net.Conn
		conexion, err = listener.Accept()
		if err != nil {
			fmt.Println("Error al aceptar conexión:")
			fmt.Println(err)
			continue
		}

		//creamos una goroutine (hilo) para manejar el cliente
		go manejarConexionCliente(conexion)
	}
}

func manejarConexionCliente(conexion net.Conn) {

	//cuando la función termine, cerramos la conexion del cliente
	defer conexion.Close()

	fmt.Println("Cliente conectado desde:")
	fmt.Println(conexion.RemoteAddr()) //obtenemos la ip del cliente

	//mantenemos la conexion abierta
	select {}
}
