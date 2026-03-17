package main

import (
	"encoding/json"
	"servidor/protocolo"
)

func GMensajePrivado(cliente *Cliente, msg *protocolo.TextMessage) {

	if !clienteIdentificado(cliente) {
		return
	}

	if msg.Username == "" || msg.Text == "" {
		GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.ResultadoInvalido)
		return
	}

	if msg.Username == cliente.nombreUsuario {
		return
	}

	clienteDestino, existe := cliente.servidor.clientes[msg.Username]

	if !existe {
		GResponderErrorExtra(cliente, protocolo.Text,
			protocolo.NoSuchUser, msg.Username)
		return
	}

	mensaje := protocolo.TextFromMessage{
		Type:     "TEXT_FROM",
		Username: cliente.nombreUsuario,
		Text:     msg.Text,
	}

	datosJSON, _ := json.Marshal(mensaje)
	clienteDestino.conn.Write(append(datosJSON, '\n'))
}

func GMensajePublico(cliente *Cliente, msg *protocolo.PublicTextFromMessage) {
	if !clienteIdentificado(cliente) {
		return
	}

	if msg.Text == "" {
		GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.ResultadoInvalido)
		return
	}

	mensaje := protocolo.PublicTextFromMessage{
		Type:     "PUBLIC_TEXT_FROM",
		Username: cliente.nombreUsuario,
		Text:     msg.Text,
	}

	GNotificarATodos(cliente.servidor, mensaje, cliente)

}
