package main

import "servidor/protocolo"

func GIdentifica(cliente *Cliente, msg *protocolo.IdentifyMessage) {
	if msg.Username == "" || len(msg.Username) > 8 {
		GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.ResultadoInvalido)
		return
	}

	if cliente.nombreUsuario != "" {
		GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.ResultadoInvalido)
		return
	}

	//map access multi-value return -> mapaccess2
	_, existe := cliente.servidor.clientes[msg.Username]
	if existe == true {
		GResponderError(cliente, protocolo.Identify, protocolo.UserAlreadyExists)
		return
	}

	cliente.nombreUsuario = msg.Username
	cliente.estado = protocolo.Active
	cliente.servidor.clientes[msg.Username] = cliente

	GResponderSuccess(cliente, protocolo.Identify)
	notificarNuevoUsuario(cliente)
}

func notificarNuevoUsuario(cliente *Cliente) {

	msg := protocolo.NewUserMessage{
		Type:     "NEW_USER",
		Username: cliente.nombreUsuario,
	}

	GNotificarATodos(cliente.servidor, msg, cliente)
}

func GActualizaStatus(cliente *Cliente, msg *protocolo.StatusMessage) {
	if !clienteIdentificado(cliente) {
		return
	}

	if (msg.Status != protocolo.Active) && (msg.Status != protocolo.Away) && (msg.Status != protocolo.Busy) {
		GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.ResultadoInvalido)
		return
	}

	cliente.estado = msg.Status
	notificarNuevoStatusDeUsuario(cliente)
	GResponderSuccess(cliente, "STATUS")
}

func notificarNuevoStatusDeUsuario(cliente *Cliente) {
	msg := protocolo.NewStatusMessage{
		Type:     "NEW_STATUS",
		Username: cliente.nombreUsuario,
		Status:   cliente.estado,
	}

	GNotificarATodos(cliente.servidor, msg, cliente)
}

func GListaDeUsuarios(cliente *Cliente) {
	if !clienteIdentificado(cliente) {
		return
	}

	listaUsuarios := make(map[string]protocolo.StatusCliente)

	for nombreUsuario, c := range cliente.servidor.clientes {
		listaUsuarios[nombreUsuario] = c.estado
	}

	datosJSON := protocolo.UserListMessage{
		Type:  "USER_LIST",
		Users: listaUsuarios,
	}

	GEnviarJSON(cliente, datosJSON)
}
