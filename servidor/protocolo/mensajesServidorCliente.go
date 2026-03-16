package protocolo

//Creamos un objeto por cada JSON que va a enviar el servidor

// ===== OPERACIONES ===== //

// respuesta del servidor a una operación como identify o enviar un mensaje
type ResponseMessage struct {
	Type      string        `json:"type"`
	Operation TipoMensaje   `json:"operation"`
	Result    TipoResultado `json:"result"`
	Extra     string        `json:"extra,omitempty"`
}

// ===== EVENTOS DEL CHAT ===== //

// mensaje que se envía al resto de clientes
// cuando un usuario se conecta e identifica
type NewUserMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
}

// mensaje que se envía cuando un usuario cambia su estado
type NewStatusMessage struct {
	Type     string        `json:"type"`
	Username string        `json:"username"`
	Status   StatusCliente `json:"status"`
}

// mensaje que se envía cuando un usuario se desconecta
type DisconnectedMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
}

// ===== LISTAS ===== //

// se envía como respuesta a USERS, se envia un diccionario con los clientes donde
// la llave es el nombre de usuario y el valor su status
type UserListMessage struct {
	Type  string                   `json:"type"`
	Users map[string]StatusCliente `json:"users"`
}

// se envía como respuesta a ROOM_USERS, se envia un diccionario con los clientes
// del cuarto donde la llave es el nombre de usuario y el valor su status
type RoomUserListMessage struct {
	Type     string                   `json:"type"`
	Roomname string                   `json:"roomname"`
	Users    map[string]StatusCliente `json:"users"`
}

// ===== MENSAJES ===== //

// mensaje privado de un usuario a otro
type TextFromMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

// mensaje público en el chat.
type PublicTextFromMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

// ==== CUARTOS ==== //

// se envía cuando un usuario se une a un cuarto
type JoinedRoomMessage struct {
	Type     string `json:"type"`
	Roomname string `json:"roomname"`
	Username string `json:"username"`
}

// se envía cuando un usuario abandona un cuarto
type LeftRoomMessage struct {
	Type     string `json:"type"`
	Roomname string `json:"roomname"`
	Username string `json:"username"`
}

// mensaje enviado dentro de un cuarto
type RoomTextFromMessage struct {
	Type     string `json:"type"`
	Roomname string `json:"roomname"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

// se envía cuando un usuario recibe una invitación a un cuarto
type InvitationMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Roomname string `json:"roomname"`
}
