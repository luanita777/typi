package main

import (
	"bufio"
	"fmt"
	"net"
	"servidor/protocolo"
)

// Tipo de dato cliente que tiene como propiedades su conexion(FD)
// y su nombre de usuario, su estado y además una referencia al servidor
// que lo instanció
type Cliente struct {
	conn          net.Conn
	servidor      *Servidor
	nombreUsuario string
	estado        protocolo.StatusCliente
}

// Creamos un nuevo cliente (Constructor)
func newCliente(conn net.Conn, servidor *Servidor) *Cliente {
	var cliente Cliente
	cliente.conn = conn
	cliente.servidor = servidor
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
		c.servidor.ProcesarMensaje(c, mensaje)
	}
}
