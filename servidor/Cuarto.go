package main

type Cuarto struct {
	nombreCuarto  string
	participantes map[string]*Cliente
}

func NuevoCuarto(nombre string) *Cuarto {
	var cuarto Cuarto
	cuarto.nombreCuarto = nombre
	cuarto.participantes = make(map[string]*Cliente)
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
