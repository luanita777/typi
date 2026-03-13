package main

import (
	"bufio"
	"fmt"
	"net"
)

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

func (c *Cliente) leeMensajesCliente() {
	var lector *bufio.Reader
	lector = bufio.NewReader(c.conn)
	for {
		var mensaje string
		var err error
		mensaje, err = lector.ReadString('\n')
		if err != nil {
			fmt.Println("El cliente se desconectó")
			c.conn.Close()
			return
		}
		fmt.Println("Mensaje del cliente: ", mensaje)
	}
}
