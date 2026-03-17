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
