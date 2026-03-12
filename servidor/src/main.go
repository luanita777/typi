package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Uso:")
		fmt.Println("go run main.go <puerto>")
		return
	}

	var puerto string = os.Args[1]

	var servidor *Servidor = newServidor(puerto)

	servidor.iniciaServidor()
}
