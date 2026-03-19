package main

import "servidor/protocolo"

func GMensajePrivado(cliente *Cliente, msg *protocolo.TextMessage) {

	if !clienteIdentificado(cliente) {
		return
	}

	if msg.Username == "" || msg.Text == "" {
		GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.ResultadoInvalido)
		return
	}

	if msg.Username == cliente.ObtenerNombreUsuario() {
		return
	}

	cliente.servidor.mutex.Lock()
	clienteDestino, existe := cliente.servidor.clientes[msg.Username]
	cliente.servidor.mutex.Unlock()

	if !existe {
		GResponderErrorExtra(cliente, protocolo.Text, protocolo.NoSuchUser, msg.Username)
		return
	}

	mensaje := protocolo.TextFromMessage{
		Type:     "TEXT_FROM",
		Username: cliente.ObtenerNombreUsuario(),
		Text:     msg.Text,
	}

	GEnviarJSON(clienteDestino, mensaje)
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
		Username: cliente.ObtenerNombreUsuario(),
		Text:     msg.Text,
	}

	GNotificarATodos(cliente.servidor, mensaje, cliente)

}
