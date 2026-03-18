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
