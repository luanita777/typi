package main

import "net"

// Tipo de dato cliente que tiene como propiedades su conexion(FD)
// y su tipo de usuario
type Cliente struct {
	conn          net.Conn
	nombreUsuario string
}

// Creamos un nuevo cliente (Constructor)
func newCliente(conn net.Conn) *Cliente {
	var cliente Cliente
	cliente.conn = conn
	return &cliente
}

//Más adelante aqui haremos cosas como leer mensajes
//del cliente y enviarle mensajes al  cliente, etc...
