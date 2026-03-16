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

func GResponderOperacionInvalida(cliente *Cliente, operacion string, respuesta string) {
	respuestaJSON := protocolo.ResponseMessage{
		Type:      "RESPONSE",
		Operation: protocolo.TipoMensaje(operacion),
		Result:    protocolo.TipoResultado(respuesta),
	}

	datosJSON, _ := json.Marshal(respuestaJSON)
	cliente.conn.Write(append(datosJSON, '\n'))
	cliente.conn.Close()
}

func GResponderError(cliente *Cliente, operacion string, respuesta string) {
	respuestaJSON := protocolo.ResponseMessage{
		Type:      "RESPONSE",
		Operation: protocolo.TipoMensaje(operacion),
		Result:    protocolo.TipoResultado(respuesta),
	}

	datosJSON, _ := json.Marshal(respuestaJSON)
	cliente.conn.Write(append(datosJSON, '\n'))
}

func GResponderSuccess(cliente *Cliente, operacion string) {
	respuesta := protocolo.ResponseMessage{
		Type:      "RESPONSE",
		Operation: protocolo.TipoMensaje(operacion),
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
		GResponderOperacionInvalida(cliente, "IDENTIFY", "NOT_IDENTIFIED")
		return false
	}

	return true
}

// ========= FUNCIONES PRINCIPALES ========= //

func GIdentifica(cliente *Cliente, msg *protocolo.IdentifyMessage) {
	if msg.Username == "" || len(msg.Username) > 8 {
		GResponderOperacionInvalida(cliente, "INVALID", "INVALID")
		return
	}

	if cliente.nombreUsuario != "" {
		GResponderOperacionInvalida(cliente, "INVALID", "INVALID")
		return
	}

	//map access multi-value return -> mapaccess2
	_, existe := cliente.servidor.clientes[msg.Username]
	if existe == true {
		GResponderError(cliente, "IDENTIFY", "USER_ALREADY_EXISTS")
		return
	}

	cliente.nombreUsuario = msg.Username
	cliente.estado = protocolo.Active
	cliente.servidor.clientes[msg.Username] = cliente

	GResponderSuccess(cliente, "IDENTIFY")
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

	if (msg.Status != "ACTIVE") && (msg.Status != "AWAY") && (msg.Status != "BUSY") {
		GResponderOperacionInvalida(cliente, "IDENTIFY", "INVALID")
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
