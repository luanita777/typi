package main

import (
	"encoding/json"
	"fmt"
	"servidor/protocolo"
)

// A lo largo de este modulo usaremos la sintaxis corta de Go para no extenderlo
// tanto. Por ejemplo:
//             Usaremos    -> jsonData, err := json.Marshal(mensaje)
//             En lugar de -> var respuestaJSON []byte
//                            var err error
//	                      respuestaJSON, err = json.Marshal(respuesta)
// Nota: La sintaxis corta es lo usual en Go, sin embargo a lo largo del
// proyecto usamos la larga por claridad.

// ========= FUNCIONES AUXILIARES ========= //

func GResponderOperacionInvalida(cliente *Cliente, operacion protocolo.TipoMensaje, respuesta protocolo.TipoResultado) {
	respuestaJSON := protocolo.ResponseMessage{
		Type:      "RESPONSE",
		Operation: operacion,
		Result:    respuesta,
	}

	datosJSON, _ := json.Marshal(respuestaJSON)
	cliente.conn.Write(append(datosJSON, '\n'))
	cliente.conn.Close()
}

func GResponderError(cliente *Cliente, operacion protocolo.TipoMensaje, respuesta protocolo.TipoResultado) {
	respuestaJSON := protocolo.ResponseMessage{
		Type:      "RESPONSE",
		Operation: operacion,
		Result:    respuesta,
	}

	datosJSON, _ := json.Marshal(respuestaJSON)
	cliente.conn.Write(append(datosJSON, '\n'))
}

func GResponderErrorExtra(cliente *Cliente, operacion protocolo.TipoMensaje, respuesta protocolo.TipoResultado, extra string) {
	respuestaJSON := protocolo.ResponseMessage{
		Type:      "RESPONSE",
		Operation: operacion,
		Result:    respuesta,
		Extra:     extra,
	}

	datosJSON, _ := json.Marshal(respuestaJSON)
	cliente.conn.Write(append(datosJSON, '\n'))
}

func GResponderSuccess(cliente *Cliente, operacion protocolo.TipoMensaje) {
	respuesta := protocolo.ResponseMessage{
		Type:      "RESPONSE",
		Operation: operacion,
		Result:    "SUCCESS",
	}

	GEnviarJSON(cliente, respuesta)
}

func GEnviarJSON(cliente *Cliente, mensaje any) {

	jsonData, err := json.Marshal(mensaje)
	if err != nil {
		fmt.Println("Error en el servidor. (Error enviando JSON)", err)
		return
	}
	cliente.conn.Write(append(jsonData, '\n'))
}

func GNotificarATodos(servidor *Servidor, mensaje any, excluir *Cliente) {

	for _, cliente := range servidor.clientes {

		if cliente == excluir {
			continue
		}

		GEnviarJSON(cliente, mensaje)
	}
}

func clienteIdentificado(cliente *Cliente) bool {

	if cliente.nombreUsuario == "" {
		GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.NotIdentified)
		return false
	}

	return true
}

// ========= FUNCIONES PRINCIPALES ========= //

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
	GNotificarNuevoUsuario(cliente)
}

func GNotificarNuevoUsuario(cliente *Cliente) {

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
	GNotificarNuevoStatusDeUsuario(cliente)
	GResponderSuccess(cliente, "STATUS")
}

func GNotificarNuevoStatusDeUsuario(cliente *Cliente) {
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
