package main

import (
	"fmt"
	"net"
)

// Creamos un tipo de dato servidor, el cuál tiene como
// propiedades un puerto, un escucha para la conexion
// y un diccionario de los clientes conectados
type Servidor struct {
	puerto   string
	listener net.Listener
	clientes map[string]*Cliente
}

// Creamos un nuevo servidor (Constructor)
func newServidor(puerto string) *Servidor {
	var servidor Servidor
	servidor.puerto = puerto
	servidor.clientes = make(map[string]*Cliente)
	return &servidor
}

// Funcion que arranca el servidor creando la direccion
// de este, luego creamos el servidor con esta dirección
// y empezamos a aceptar clientes
func (s *Servidor) iniciaServidor() {

	var direccion string = ":" + s.puerto
	var err error
	s.listener, err = net.Listen("tcp", direccion)

	if err != nil {
		fmt.Println("Error iniciando servidor")
		panic(err)
	}

	fmt.Println("Servidor escuchando en puerto", s.puerto)
	for {
		var conexion net.Conn
		conexion, err = s.listener.Accept()
		if err != nil {
			fmt.Println("Error aceptando conexión")
			continue
		}
		go s.manejarConexionCliente(conexion)
	}
}

// Creamos un hilo que va a manejar al cliente
func (s *Servidor) manejarConexionCliente(conn net.Conn) {

	var cliente *Cliente = newCliente(conn)
	fmt.Println("Cliente conectado:", conn.RemoteAddr())

	//descartamos cliente porque todavia no lo usamos
	_ = cliente
}
