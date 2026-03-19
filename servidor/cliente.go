package main

import (
	"bufio"
	"fmt"
	"net"
	"servidor/protocolo"
	"sync"
)

// Tipo de dato cliente que tiene como propiedades su conexion(FD)
// y su nombre de usuario, su estado y además una referencia al servidor
// que lo instanció
type Cliente struct {
	conn          net.Conn
	servidor      *Servidor
	nombreUsuario string
	estado        protocolo.StatusCliente
	mutex         sync.Mutex
}

// Creamos un nuevo cliente (Constructor)
func newCliente(conn net.Conn, servidor *Servidor) *Cliente {
	var cliente Cliente
	cliente.conn = conn
	cliente.servidor = servidor
	return &cliente
}

func (c *Cliente) ObtenerNombreUsuario() string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.nombreUsuario
}

func (c *Cliente) ObtenerEstado() protocolo.StatusCliente {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.estado
}

func (c *Cliente) leeMensajesCliente() {

	//nos aseguramos de siempre cerrar la conexion al terminar la funcion
	defer func() {
		fmt.Println("Limpiando socket")
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
			nombre := c.ObtenerNombreUsuario()
			if nombre == "" {
				fmt.Println("Deteniendo lectura de mensajes para cliente no identificado.")
			} else {
				fmt.Printf("Deteniendo lectura de mensajes para %s por error de protocolo.\n", c.ObtenerNombreUsuario())
			}
			return
		}
	}
}
