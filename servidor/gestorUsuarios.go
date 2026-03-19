package main

import (
	"servidor/protocolo"
)

func GIdentifica(cliente *Cliente, msg *protocolo.IdentifyMessage) {
	if msg.Username == "" || len(msg.Username) > 8 {
		GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.ResultadoInvalido)
		return
	}

	if cliente.ObtenerNombreUsuario() != "" {
		GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.ResultadoInvalido)
		return
	}

	//map access multi-value return -> mapaccess2
	cliente.servidor.mutex.Lock()
	_, existe := cliente.servidor.clientes[msg.Username]
	if existe {
		cliente.servidor.mutex.Unlock()
		GResponderError(cliente, protocolo.Identify, protocolo.UserAlreadyExists)
		return
	}

	cliente.nombreUsuario = msg.Username
	cliente.estado = protocolo.Active
	cliente.servidor.clientes[msg.Username] = cliente
	cliente.servidor.mutex.Unlock()

	GResponderSuccess(cliente, protocolo.Identify)

	mensajeJSON := protocolo.NewUserMessage{
		Type:     "NEW_USER",
		Username: cliente.ObtenerNombreUsuario(),
	}

	GNotificarATodos(cliente.servidor, mensajeJSON, cliente)
}

func GActualizaStatus(cliente *Cliente, msg *protocolo.StatusMessage) {
	if !clienteIdentificado(cliente) {
		return
	}

	if (msg.Status != protocolo.Active) && (msg.Status != protocolo.Away) && (msg.Status != protocolo.Busy) {
		GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.ResultadoInvalido)
		return
	}

	cliente.mutex.Lock()
	cliente.estado = msg.Status
	cliente.mutex.Unlock()

	mensajeJSON := protocolo.NewStatusMessage{
		Type:     "NEW_STATUS",
		Username: cliente.ObtenerNombreUsuario(),
		Status:   cliente.ObtenerEstado(),
	}
	GNotificarATodos(cliente.servidor, mensajeJSON, cliente)

	GResponderSuccess(cliente, protocolo.Status)
}

func GListaDeUsuarios(cliente *Cliente) {
	if !clienteIdentificado(cliente) {
		return
	}

	listaUsuarios := make(map[string]protocolo.StatusCliente)

	cliente.servidor.mutex.Lock()
	copiaClientes := make([]*Cliente, 0, len(cliente.servidor.clientes))
	for _, c := range cliente.servidor.clientes {
		copiaClientes = append(copiaClientes, c)
	}
	cliente.servidor.mutex.Unlock()

	for _, c := range copiaClientes {
		listaUsuarios[c.ObtenerNombreUsuario()] = c.ObtenerEstado()
	}

	datosJSON := protocolo.UserListMessage{
		Type:  "USER_LIST",
		Users: listaUsuarios,
	}

	GEnviarJSON(cliente, datosJSON)
}

func GDesconecta(cliente *Cliente, msg *protocolo.DisconnectMessage) {

	nombre := cliente.ObtenerNombreUsuario()

	cliente.servidor.mutex.Lock()

	if nombre == "" {
		cliente.servidor.mutex.Unlock()
		cliente.conn.Close()
		return
	}

	_, existe := cliente.servidor.clientes[nombre]
	if !existe {
		cliente.servidor.mutex.Unlock()
		cliente.conn.Close()
		return
	}

	//copiamos la lista de cuartos
	copiaCuartos := make([]*Cuarto, 0, len(cliente.servidor.cuartos))
	for _, cuarto := range cliente.servidor.cuartos {
		copiaCuartos = append(copiaCuartos, cuarto)
	}

	delete(cliente.servidor.clientes, nombre)
	cliente.servidor.mutex.Unlock()

	for _, cuarto := range copiaCuartos {
		if cuarto.EstaEnCuarto(nombre) {
			msg := protocolo.LeaveRoomMessage{
				Type:     protocolo.LeaveRoom,
				Roomname: cuarto.nombreCuarto,
			}

			GAbandonarCuarto(cliente, &msg)
		}
	}

	mensajeJSON := protocolo.DisconnectedMessage{
		Type:     "DISCONNECTED",
		Username: nombre,
	}

	GNotificarATodos(cliente.servidor, mensajeJSON, cliente)
	cliente.conn.Close()
}
