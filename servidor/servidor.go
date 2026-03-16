package main

import (
	"encoding/json"
	"fmt"
	"net"
	"servidor/protocolo"
)

// Creamos un tipo de dato servidor, el cuál tiene como
// propiedades un puerto, un escucha para la conexion
// y un diccionario de los clientes conectados
type Servidor struct {
	puerto      string
	listener    net.Listener
	clientes    map[string]*Cliente
	numClientes int
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

	var cliente *Cliente = newCliente(conn, s)
	s.numClientes++
	fmt.Printf("Cliente %v conectado.\n", s.clientes)
	go cliente.leeMensajesCliente()
}

// En esencia lo que se hace es que se recibe un mensaje dado un cliente
// y se convierte ese mensaje a bytes para poder operarlo con Unmarshal
// esto decodifica el JSON y lo pasa a un struct; luego, usamos el struct
// que creamos que unicamente tiene el campo Type, usamos esto para
// obtener el tipo que nos mandaron en el mensaje, y ya sabiendo el tipo
// entonces podemos decidir qué debemos hacer de acuerdo al mensaje enviado
// por ejemplo, si nos mandaron un identify debemos revisar si existe
// ya ese cliente, si sí respondemos que no puede usar ese nombre y si no
// le decimos que sucess y lo dejamos entrar al chat
func (s *Servidor) ProcesarMensaje(cliente *Cliente, mensaje string) {

	var mensajeJSON []byte = []byte(mensaje)
	var mensajeBase protocolo.MensajeBase
	var err error
	err = json.Unmarshal(mensajeJSON, &mensajeBase)
	if err != nil {
		fmt.Println("JSON inválido: ", err)
		return
	}

	switch mensajeBase.Type {
	case protocolo.Identify:
		var msj protocolo.IdentifyMessage
		var err error = json.Unmarshal(mensajeJSON, &msj)
		if err != nil {
			GResponderMensajeInvalido(cliente)
			return
		}
		GIdentifica(cliente, &msj)

	case protocolo.Status:
		var msj protocolo.StatusMessage
		json.Unmarshal(mensajeJSON, &msj)
		//gestor.gStatus(cliente, &msj)

	}

}
