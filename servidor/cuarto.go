package main

import "sync"

// aqui lo idel sería tener algo así como un conjunto para los invitados
// pero como eso no existe en Go lo simulamos con un diccionario de strings y
// booleanos, así agregamos y eliminamso en O(1) y nos facilitamos la vida con la sintaxis
type Cuarto struct {
	nombreCuarto  string
	participantes map[string]*Cliente
	invitados     map[string]bool
	mutex         sync.Mutex
}

func NuevoCuarto(nombre string) *Cuarto {
	var cuarto Cuarto
	cuarto.nombreCuarto = nombre
	cuarto.participantes = make(map[string]*Cliente)
	cuarto.invitados = make(map[string]bool)
	return &cuarto
}

func (c *Cuarto) AgregarCliente(cliente *Cliente) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.participantes[cliente.ObtenerNombreUsuario()] = cliente
}

func (c *Cuarto) EliminarCliente(cliente *Cliente) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.participantes, cliente.ObtenerNombreUsuario())
}

func (c *Cuarto) EstaEnCuarto(nombreUsuario string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, existe := c.participantes[nombreUsuario]
	return existe
}

func (c *Cuarto) AgregarInvitado(cliente *Cliente) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.invitados[cliente.ObtenerNombreUsuario()] = true
}

func (c *Cuarto) EliminarInvitado(cliente *Cliente) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.invitados, cliente.ObtenerNombreUsuario())
}

func (c *Cuarto) EstaInvitado(nombreUsuario string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, existe := c.invitados[nombreUsuario]
	return existe
}

func (c *Cuarto) UnirCliente(cliente *Cliente) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	nombre := cliente.ObtenerNombreUsuario()
	c.participantes[nombre] = cliente
	delete(c.invitados, nombre)
}

func (c *Cuarto) NumeroParticipantes() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return len(c.participantes)
}

// Aquí usamos mutex mientras copiamos nuestros clientes a una nueva
// estructura y luego desbloqueamos,  hacemos esto para poder usar
// esa lista de copia de nuestros participantes en el cuarto para
// poder enviarles mensaje sin bloquear el cuarto mientras los enviamos
func (c *Cuarto) ObtenerParticipantes() []*Cliente {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	participantes := make([]*Cliente, 0, len(c.participantes))
	for _, participante := range c.participantes {
		participantes = append(participantes, participante)
	}

	return participantes
}

func (c *Cuarto) EnviarMensajeCuarto(mensaje any, excluir *Cliente) {
	participantesCuarto := c.ObtenerParticipantes()

	for _, cliente := range participantesCuarto {
		if cliente == excluir {
			continue
		}
		GEnviarJSON(cliente, mensaje)
	}
}
