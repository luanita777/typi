package protocol

//Creamos un objeti por cada JSON que recibe el servidor, así
// no creamos manualmente los JSON cuando los queramos enviar
// Estos "objetos" se crean con structs que se comportan tal que:
//
// Sintaxis general:
//
//     type NombreDelStruct struct {
//         Campo Tipo `tag`
//     }
//
// - `type` define un nuevo tipo.
// - `NombreDelStruct` es el nombre del tipo que estamos creando.
// - `struct` indica que el tipo es una estructura de datos.
// - Dentro de las llaves `{}` se definen los campos del struct.
//
// Donde cada campo de nuestro objeto tiene:
//
//     NombreDelCampo TipoDelCampo `tag`
//
// Ejemplo:
//
//     Username string `json:"username"`
//
// - `Username` es el nombre del campo en Go.
// - `string` es el tipo del campo.
// - `json:"username"`: es un struct tag que indica cómo se llama
//   el campo cuando se convierte a JSON.
//
// Nota: Los nombres de los campos comienzan con mayúscula porque en Go eso
// significa que son exportados (visibles fuera del paquete). Esto es
// necesario para que la biblioteca de JSON pueda leer y escribir
// esos campos al convertir entre structs y JSON.

// mensaje IDENTIFY enviado por el cliente para identificarse
type IdentifyMessage struct {
	Type     TipoMensaje `json:"type"`
	Username string      `json:"username"`
}

// mensaje STATUS enviado por el cliente para notificar del cambio de estado
type StatusMessage struct {
	Type   TipoMensaje   `json:"type"`
	Status StatusCliente `json:"status"`
}

// mensaje USERS que pide la lista de usuarios en el chat
type UsersMessage struct {
	Type TipoMensaje `json:"type"`
}

// mensaje privado TEXT de un cliente a otro
type TextMessage struct {
	Type     TipoMensaje `json:"type"`
	Username string      `json:"username"`
	Text     string      `json:"text"`
}

// mensaje público PUBLIC_TEXT en el chat grupal
type PublicTextMessage struct {
	Type TipoMensaje `json:"type"`
	Text string      `json:"text"`
}

// mensaje NEW_ROOM para crear un nuevo cuarto
type NewRoomMessage struct {
	Type     TipoMensaje `json:"type"`
	Roomname string      `json:"roomname"`
}

// mensaje de invitacion INVITE para unirse a un cuarto
type InviteMessage struct {
	Type      TipoMensaje `json:"type"`
	Roomname  string      `json:"roomname"`
	Usernames []string    `json:"usernames"`
}

// mensaje JOIN_ROOM que permite unirse a un cuarto
type JoinRoomMessage struct {
	Type     TipoMensaje `json:"type"`
	Roomname string      `json:"roomname"`
}

// mensaje ROOM_USERS que regresa la lista de usuarios en el cuarto
type RoomUsersMessage struct {
	Type     TipoMensaje `json:"type"`
	Roomname string      `json:"roomname"`
}

// mensaje ROOM_TEXT para mandar un mensaje al cuarto
type RoomTextMessage struct {
	Type     TipoMensaje `json:"type"`
	Roomname string      `json:"roomname"`
	Text     string      `json:"text"`
}

// mensaje LEAVE_ROOM para abandonar el cuarto
type LeaveRoomMessage struct {
	Type     TipoMensaje `json:"type"`
	Roomname string      `json:"roomname"`
}

// mensaje DISCONNECT para salirse de la aplicacion
type DisconnectMessage struct {
	Type TipoMensaje `json:"type"`
}
