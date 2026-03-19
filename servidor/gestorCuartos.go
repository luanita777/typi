package main

import (
	//	"encoding/json"
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

	_, existe := cliente.servidor.cuartos[msg.Roomname]
	if existe == true {
		GResponderErrorExtra(cliente, protocolo.NewRoom, protocolo.RoomAlreadyExists, msg.Roomname)
		return
	}

	cuarto := NuevoCuarto(msg.Roomname)
	cliente.servidor.cuartos[msg.Roomname] = cuarto
	cuarto.AgregarCliente(cliente)

	GResponderSuccessExtra(cliente, protocolo.NewRoom, msg.Roomname)
}

func GInvitaACuarto(cliente *Cliente, msg *protocolo.InviteMessage) {

	if !clienteIdentificado(cliente) {
		return
	}

	cuartoActual, existeCuarto := cliente.servidor.cuartos[msg.Roomname]
	if !existeCuarto {
		GResponderErrorExtra(cliente, protocolo.Invite, protocolo.NoSuchRoom, msg.Roomname)
		return
	}

	_, estaEnElCuarto := cuartoActual.participantes[cliente.nombreUsuario]
	if !estaEnElCuarto {
		GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.NotJoined)
		return
	}

	if msg.Usernames == nil || len(msg.Usernames) == 0 {
		GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.ResultadoInvalido)
		return
	}

	for _, username := range msg.Usernames {
		if username == "" {
			GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.ResultadoInvalido)
			return
		}

		invitado, existe := cliente.servidor.clientes[username]
		if !existe {
			GResponderErrorExtra(cliente, protocolo.Invite, protocolo.NoSuchUser, username)
			return
		}

		if cuartoActual.EstaInvitado(invitado.nombreUsuario) ||
			cuartoActual.EstaEnCuarto(invitado.nombreUsuario) {
			continue
		}

		mensajeJSON := protocolo.InvitationMessage{
			Type:     "INVITATION",
			Username: cliente.nombreUsuario,
			Roomname: msg.Roomname,
		}

		cuartoActual.invitados[invitado.nombreUsuario] = true

		GEnviarJSON(invitado, mensajeJSON)
	}
}

func GUnirseACuarto(cliente *Cliente, msg *protocolo.JoinRoomMessage) {
	if !clienteIdentificado(cliente) {
		return
	}

	cuartoActual, existeCuarto := cliente.servidor.cuartos[msg.Roomname]
	if !existeCuarto {
		GResponderErrorExtra(cliente, protocolo.JoinRoom, protocolo.NoSuchRoom, msg.Roomname)
		return
	}

	if !cuartoActual.EstaInvitado(cliente.nombreUsuario) {
		GResponderErrorExtra(cliente, protocolo.JoinRoom, protocolo.NotInvited, msg.Roomname)
		return
	}

	if cuartoActual.EstaEnCuarto(cliente.nombreUsuario) {
		return
	}

	cuartoActual.AgregarCliente(cliente)
	delete(cuartoActual.invitados, cliente.nombreUsuario)

	GResponderSuccessExtra(cliente, protocolo.JoinRoom, msg.Roomname)
	notificarUnionACuartoATodos(cliente, cuartoActual)

}

func notificarUnionACuartoATodos(cliente *Cliente, cuarto *Cuarto) {
	mensajeJSON := protocolo.JoinedRoomMessage{
		Type:     "JOINED_ROOM",
		Roomname: cuarto.nombreCuarto,
		Username: cliente.nombreUsuario,
	}
	enviarATodosEnElCuarto(cuarto, cliente, mensajeJSON)
}

func enviarATodosEnElCuarto(cuarto *Cuarto, excluir *Cliente, mensaje any) {
	for _, cliente := range cuarto.participantes {
		if cliente == excluir {
			continue
		}
		GEnviarJSON(cliente, mensaje)
	}
}

func GUsuariosCuarto(cliente *Cliente, msg *protocolo.RoomUsersMessage) {
	if !clienteIdentificado(cliente) {
		return
	}

	cuartoActual, existeCuarto := cliente.servidor.cuartos[msg.Roomname]
	if !existeCuarto {
		GResponderErrorExtra(cliente, protocolo.RoomUsers, protocolo.NoSuchRoom, msg.Roomname)
		return
	}

	if !cuartoActual.EstaEnCuarto(cliente.nombreUsuario) {
		GResponderErrorExtra(cliente, protocolo.RoomUsers, protocolo.NotJoined, msg.Roomname)
		return
	}

	listaUsuarios := make(map[string]protocolo.StatusCliente)
	for nombreUsuario, c := range cuartoActual.participantes {
		listaUsuarios[nombreUsuario] = c.estado
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

	cuartoActual, existeCuarto := cliente.servidor.cuartos[msg.Roomname]
	if !existeCuarto {
		GResponderErrorExtra(cliente, protocolo.RoomText, protocolo.NoSuchRoom, msg.Roomname)
		return
	}

	if !cuartoActual.EstaEnCuarto(cliente.nombreUsuario) {
		GResponderErrorExtra(cliente, protocolo.RoomText, protocolo.NotJoined, msg.Roomname)
		return
	}

	mensajeJSON := protocolo.RoomTextFromMessage{
		Type:     "ROOM_TEXT_FROM",
		Roomname: msg.Roomname,
		Username: cliente.nombreUsuario,
		Text:     msg.Text,
	}

	enviarATodosEnElCuarto(cuartoActual, cliente, mensajeJSON)
}
