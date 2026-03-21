package main

import (
	"encoding/json"
	"fmt"
	"io"
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

	GEnviarJSON(cliente, respuestaJSON)
}

func GResponderError(cliente *Cliente, operacion protocolo.TipoMensaje, respuesta protocolo.TipoResultado) {
	respuestaJSON := protocolo.ResponseMessage{
		Type:      "RESPONSE",
		Operation: operacion,
		Result:    respuesta,
	}

	GEnviarJSON(cliente, respuestaJSON)
}

func GResponderErrorExtra(cliente *Cliente, operacion protocolo.TipoMensaje, respuesta protocolo.TipoResultado, extra string) {
	respuestaJSON := protocolo.ResponseMessage{
		Type:      "RESPONSE",
		Operation: operacion,
		Result:    respuesta,
		Extra:     extra,
	}

	GEnviarJSON(cliente, respuestaJSON)
}

func GResponderSuccess(cliente *Cliente, operacion protocolo.TipoMensaje) {
	respuesta := protocolo.ResponseMessage{
		Type:      "RESPONSE",
		Operation: operacion,
		Result:    protocolo.Success,
	}

	GEnviarJSON(cliente, respuesta)
}

func GResponderSuccessExtra(cliente *Cliente, operacion protocolo.TipoMensaje, extra string) {
	respuestaJSON := protocolo.ResponseMessage{
		Type:      "RESPONSE",
		Operation: operacion,
		Result:    protocolo.Success,
		Extra:     extra,
	}

	GEnviarJSON(cliente, respuestaJSON)
}

func GEnviarJSON(cliente *Cliente, mensaje any) {

	jsonData, err := json.Marshal(mensaje)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Cliente cerró la conexion.")
			return
		} else {
			fmt.Println("Error de lectura:", err)
		}
		return
	}
	fmt.Println(string(jsonData))
	cliente.mutex.Lock()
	defer cliente.mutex.Unlock()
	data := append(jsonData, 0)
	_, err = cliente.conn.Write(append(data))
	if err != nil {
		fmt.Println("Error escribiendo en conexion: ", err)
	}

}

func GNotificarATodos(servidor *Servidor, mensaje any, excluir *Cliente) {
	servidor.mutex.Lock()
	copiaClientes := make([]*Cliente, 0, len(servidor.clientes))
	for _, c := range servidor.clientes {
		copiaClientes = append(copiaClientes, c)
	}
	servidor.mutex.Unlock()

	for _, cliente := range copiaClientes {
		if cliente == excluir {
			continue
		}

		GEnviarJSON(cliente, mensaje)
	}
}

func clienteIdentificado(cliente *Cliente) bool {

	if cliente.ObtenerNombreUsuario() == "" {
		GResponderOperacionInvalida(cliente, protocolo.Invalid, protocolo.NotIdentified)
		return false
	}

	return true
}
