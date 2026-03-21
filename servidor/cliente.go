package main

import (
	"fmt"
	"io"
	"net"
	"servidor/protocolo"
	"strings"
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

// Leemos continuamente los bytes del socket, respetando que el protocolo usa
// mensajes terminados en \0 (null terminated como dijo Canek)
//
// Como TCP no respeta los límites de los mensajes y solo manda bytes,
// si se conectan muchos a la vez y mandan muchos mensajes a la vez
// pueden pasar que un mensaje llegue incompleto, en varias partes o
// que varios mensajes lleguen juntos
//
// Entonces lo que hacemos es:
//
//  1. Leer bloques de bytes del socket con Read
//  2. Ir acumulando esos bytes en un string llamado "acumulado"
//     (esto sirve para juntar fragmentos de mensajes)
//  3. Después de cada lectura, revisamos si en "acumulado" ya existe un \0
//     que indica el final de un mensaje
//  4. Si encontramos un \0:
//     - tomamos todo lo que está ANTES de ese \0 entonces  ese es un mensaje completo
//     - lo procesamos
//     - eliminamos esa parte del acumulado (dejando lo que sobra)
//     Esto se hace en un loop interno porque pueden venir varios mensajes distintos juntos
//  5. Si en algun mensaje aparece \n automaticamente cerramos la conexion porque es invalido
//     con el protocolo

func (c *Cliente) leeMensajesCliente() {

	defer func() {
		fmt.Println("Limpiando socket")
		msg := protocolo.DisconnectMessage{
			Type: protocolo.Disconnect,
		}
		GDesconecta(c, &msg)
	}()

	buffer := make([]byte, 1024)
	acumulado := ""

	for {
		n, err := c.conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Cliente cerró la conexión")
			} else {
				fmt.Println("Error de lectura:", err)
			}
			return
		}

		// Agregamos lo leído al acumulador
		acumulado += string(buffer[:n])

		if strings.Contains(acumulado, "\n") && !strings.Contains(acumulado, "\x00") {
			fmt.Println("Cliente invalido: usa \\n en lugar de \\0. Desconectando...")
			return
		}

		// Procesamos todos los mensajes completos (\0)
		for strings.Contains(acumulado, "\x00") {

			partes := strings.SplitN(acumulado, "\x00", 2)

			mensaje := partes[0]
			acumulado = partes[1]

			fmt.Println("Procesando:", mensaje)

			continua := c.servidor.ProcesarMensaje(c, mensaje)
			if !continua {
				nombre := c.ObtenerNombreUsuario()
				if nombre == "" {
					fmt.Println("Deteniendo lectura de mensajes para cliente no identificado.")
				} else {
					fmt.Printf("Deteniendo lectura de mensajes para %s por error de protocolo.\n", nombre)
				}
				return
			}
		}
	}
}
