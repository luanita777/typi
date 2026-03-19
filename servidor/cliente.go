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

	//nos aseguramos de siempre cerrar la conexion al terminar la funcion
	defer func() {
		fmt.Println("Ejecutando limpieza del socket...")
		msg := protocolo.DisconnectMessage{
			Type: protocolo.Disconnect,
		}
		GDesconecta(c, &msg)
	}()

	lector := bufio.NewReader(c.conn)
	for {
		mensaje, err := lector.ReadString('\n')
		if err != nil {
			fmt.Println("Error de lectura o cliente desconectado.")
			return
		}

		continua := c.servidor.ProcesarMensaje(c, mensaje)
		if !continua {
			fmt.Println("Deteniendo lectura de mensajes para este cliente por error de protocolo.")
			return
		}
	}
}
