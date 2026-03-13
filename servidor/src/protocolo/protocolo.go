package protocol

// tipo de operación que el cliente solicita al servidor
type TipoMensaje string

const (
	Identify   TipoMensaje = "IDENTIFY"
	Status     TipoMensaje = "STATUS"
	Users      TipoMensaje = "USERS"
	Text       TipoMensaje = "TEXT"
	PublicText TipoMensaje = "PUBLIC_TEXT"

	NewRoom   TipoMensaje = "NEW_ROOM"
	Invite    TipoMensaje = "INVITE"
	JoinRoom  TipoMensaje = "JOIN_ROOM"
	RoomUsers TipoMensaje = "ROOM_USERS"
	RoomText  TipoMensaje = "ROOM_TEXT"
	LeaveRoom TipoMensaje = "LEAVE_ROOM"

	Disconnect TipoMensaje = "DISCONNECT"
)

// estado actual de un usuario dentro del chat
type StatusCliente string

const (
	Active StatusCliente = "ACTIVE"
	Away   StatusCliente = "AWAY"
	Busy   StatusCliente = "BUSY"
)

// resultado de alguna operación solicitada por un cliente al servidor.
type TipoResultado string

const (
	Success           TipoResultado = "SUCCESS"
	UserAlreadyExists TipoResultado = "USER_ALREADY_EXISTS"
	NoSuchUser        TipoResultado = "NO_SUCH_USER"
	NoSuchRoom        TipoResultado = "NO_SUCH_ROOM"
	NotInvited        TipoResultado = "NOT_INVITED"
	NotJoined         TipoResultado = "NOT_JOINED"
	Invalid           TipoResultado = "INVALID"
	NotIdentified     TipoResultado = "NOT_IDENTIFIED"
)
