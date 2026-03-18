package main

// aqui lo idel sería tener algo así como un conjunto para los invitados
// pero como eso no existe en Go lo simulamos con un diccionario de strings y
// booleanos, así agregamos y eliminamso en O(1) y nos facilitamos la vida con la sintaxis
type Cuarto struct {
	nombreCuarto  string
	participantes map[string]*Cliente
	invitados     map[string]bool
}

func NuevoCuarto(nombre string) *Cuarto {
	var cuarto Cuarto
	cuarto.nombreCuarto = nombre
	cuarto.participantes = make(map[string]*Cliente)
	cuarto.invitados = make(map[string]bool)
	return &cuarto
}

func (c *Cuarto) AgregarCliente(cliente *Cliente) {
	c.participantes[cliente.nombreUsuario] = cliente
}

func (c *Cuarto) EliminarCliente(cliente *Cliente) {
	delete(c.participantes, cliente.nombreUsuario)
}

func (c *Cuarto) EnviarMensajeCuarto(mensaje any, excluir *Cliente) {
	for _, cliente := range c.participantes {
		if cliente == excluir {
			continue
		}
		GEnviarJSON(cliente, mensaje)
	}
}

func (c *Cuarto) EstaEnCuarto(nombreUsuario string) bool {
	_, existe := c.participantes[nombreUsuario]
	return existe
}

func (c *Cuarto) EstaInvitado(nombreUsuario string) bool {
	_, existe := c.invitados[nombreUsuario]
	return existe
}
