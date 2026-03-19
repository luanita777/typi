package main

import (
	"servidor/protocolo"
)

func GCreaNuevoCuarto(cliente *Cliente, msg *protocolo.NewRoomMessage) {
	if !clienteIdentificado(cliente) {
		return
	}

	if msg.Roomname == "" || len(msg.Roomname) > 16 {
		GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.ResultadoInvalido)
		return
	}

	cliente.servidor.mutex.Lock()

	_, existe := cliente.servidor.cuartos[msg.Roomname]
	if existe == true {
		cliente.servidor.mutex.Unlock()
		GResponderErrorExtra(cliente, protocolo.NewRoom, protocolo.RoomAlreadyExists, msg.Roomname)
		return
	}

	cuarto := NuevoCuarto(msg.Roomname)
	cliente.servidor.cuartos[msg.Roomname] = cuarto

	cliente.servidor.mutex.Unlock()

	cuarto.AgregarCliente(cliente)
	GResponderSuccessExtra(cliente, protocolo.NewRoom, msg.Roomname)
}

func GInvitaACuarto(cliente *Cliente, msg *protocolo.InviteMessage) {

	if !clienteIdentificado(cliente) {
		return
	}

	cliente.servidor.mutex.Lock()
	cuartoActual, existeCuarto := cliente.servidor.cuartos[msg.Roomname]
	cliente.servidor.mutex.Unlock()

	if !existeCuarto {
		GResponderErrorExtra(cliente, protocolo.Invite, protocolo.NoSuchRoom, msg.Roomname)
		return
	}

	if !cuartoActual.EstaEnCuarto(cliente.ObtenerNombreUsuario()) {
		GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.NotJoined)
		return
	}

	if len(msg.Usernames) == 0 {
		GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.ResultadoInvalido)
		return
	}

	for _, username := range msg.Usernames {
		if username == "" {
			GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.ResultadoInvalido)
			return
		}

		cliente.servidor.mutex.Lock()
		invitado, existe := cliente.servidor.clientes[username]
		cliente.servidor.mutex.Unlock()

		if !existe {
			GResponderErrorExtra(cliente, protocolo.Invite, protocolo.NoSuchUser, username)
			return
		}

		if cuartoActual.EstaInvitado(invitado.ObtenerNombreUsuario()) ||
			cuartoActual.EstaEnCuarto(invitado.ObtenerNombreUsuario()) {
			continue
		}
		cuartoActual.AgregarInvitado(invitado)

		mensajeJSON := protocolo.InvitationMessage{
			Type:     "INVITATION",
			Username: cliente.ObtenerNombreUsuario(),
			Roomname: msg.Roomname,
		}

		GEnviarJSON(invitado, mensajeJSON)
	}
}

func GUnirseACuarto(cliente *Cliente, msg *protocolo.JoinRoomMessage) {
	if !clienteIdentificado(cliente) {
		return
	}

	cliente.servidor.mutex.Lock()
	cuartoActual, existeCuarto := cliente.servidor.cuartos[msg.Roomname]
	cliente.servidor.mutex.Unlock()
	if !existeCuarto {
		GResponderErrorExtra(cliente, protocolo.JoinRoom, protocolo.NoSuchRoom, msg.Roomname)
		return
	}

	if !cuartoActual.EstaInvitado(cliente.ObtenerNombreUsuario()) {
		GResponderErrorExtra(cliente, protocolo.JoinRoom, protocolo.NotInvited, msg.Roomname)
		return
	}

	if cuartoActual.EstaEnCuarto(cliente.ObtenerNombreUsuario()) {
		return
	}

	cuartoActual.UnirCliente(cliente)
	GResponderSuccessExtra(cliente, protocolo.JoinRoom, msg.Roomname)

	mensajeJSON := protocolo.JoinedRoomMessage{
		Type:     "JOINED_ROOM",
		Roomname: msg.Roomname,
		Username: cliente.ObtenerNombreUsuario(),
	}

	cuartoActual.EnviarMensajeCuarto(mensajeJSON, cliente)

}

func GUsuariosCuarto(cliente *Cliente, msg *protocolo.RoomUsersMessage) {
	if !clienteIdentificado(cliente) {
		return
	}

	cliente.servidor.mutex.Lock()
	cuartoActual, existeCuarto := cliente.servidor.cuartos[msg.Roomname]
	cliente.servidor.mutex.Unlock()

	if !existeCuarto {
		GResponderErrorExtra(cliente, protocolo.RoomUsers, protocolo.NoSuchRoom, msg.Roomname)
		return
	}

	if !cuartoActual.EstaEnCuarto(cliente.nombreUsuario) {
		GResponderErrorExtra(cliente, protocolo.RoomUsers, protocolo.NotJoined, msg.Roomname)
		return
	}

	listaUsuarios := make(map[string]protocolo.StatusCliente)

	for _, c := range cuartoActual.ObtenerParticipantes() {
		listaUsuarios[c.ObtenerNombreUsuario()] = c.ObtenerEstado()
	}

	mensajeJSON := protocolo.RoomUserListMessage{
		Type:     "ROOM_USER_LIST",
		Roomname: msg.Roomname,
		Users:    listaUsuarios,
	}

	GEnviarJSON(cliente, mensajeJSON)

}

func GRoomText(cliente *Cliente, msg *protocolo.RoomTextMessage) {
	if !clienteIdentificado(cliente) {
		return
	}

	cliente.servidor.mutex.Lock()
	cuartoActual, existeCuarto := cliente.servidor.cuartos[msg.Roomname]
	cliente.servidor.mutex.Unlock()

	if !existeCuarto {
		GResponderErrorExtra(cliente, protocolo.RoomText, protocolo.NoSuchRoom, msg.Roomname)
		return
	}

	if !cuartoActual.EstaEnCuarto(cliente.ObtenerNombreUsuario()) {
		GResponderErrorExtra(cliente, protocolo.RoomText, protocolo.NotJoined, msg.Roomname)
		return
	}

	mensajeJSON := protocolo.RoomTextFromMessage{
		Type:     "ROOM_TEXT_FROM",
		Roomname: msg.Roomname,
		Username: cliente.ObtenerNombreUsuario(),
		Text:     msg.Text,
	}

	cuartoActual.EnviarMensajeCuarto(mensajeJSON, cliente)
}

func GAbandonarCuarto(cliente *Cliente, msg *protocolo.LeaveRoomMessage) {
	if !clienteIdentificado(cliente) {
		return
	}

	cliente.servidor.mutex.Lock()
	cuartoActual, existeCuarto := cliente.servidor.cuartos[msg.Roomname]
	cliente.servidor.mutex.Unlock()

	if !existeCuarto {
		GResponderErrorExtra(cliente, protocolo.LeaveRoom, protocolo.NoSuchRoom, msg.Roomname)
		return
	}

	if !cuartoActual.EstaEnCuarto(cliente.ObtenerNombreUsuario()) {
		GResponderErrorExtra(cliente, protocolo.LeaveRoom, protocolo.NotJoined, msg.Roomname)
		return
	}

	cuartoActual.EliminarCliente(cliente)

	mensajeJSON := protocolo.LeftRoomMessage{
		Type:     "LEFT_ROOM",
		Roomname: msg.Roomname,
		Username: cliente.ObtenerNombreUsuario(),
	}

	cuartoActual.EnviarMensajeCuarto(mensajeJSON, cliente)

	if cuartoActual.NumeroParticipantes() == 0 {
		cliente.servidor.mutex.Lock()
		delete(cliente.servidor.cuartos, cuartoActual.nombreCuarto)
		cliente.servidor.mutex.Unlock()
	}
}
