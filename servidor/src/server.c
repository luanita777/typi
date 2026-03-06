#include <stdio.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>

// FUNCION AUXILIAR //
// Creamos un nuevo socket especificando que usaremos el protocolo IPv4 (AF_INET)
// y especificamos que usaremos comunicacion TCP pues queremos que todo llegue en
// orden (SOCK_STREAM)
int createSocket(){
  return socket(AF_INET, SOCK_STREAM, 0);
}

// FUNCION AUXILIAR //
// Crea una dirección para nuestro socket
struct sockaddr_in* createAddress(char *ip, int port){

  // Asignamos memoria para guardar nuestra dirección
  struct sockaddr_in *address = malloc(sizeof(struct sockaddr_in));

  // Limpiamos la memoria que usaremos asignandola toda a cero
  memset(address, 0, sizeof(struct sockaddr_in));

  // Definimos el protocolo IPv4 y convertimos el puerto a un formato especial
  // para la red; usamos -> porque address es un apuntador
  address->sin_family = AF_INET;
  address->sin_port = htons(port);

  // Convierte nuestra IP en binario para que el OS lo lea instantaneamente, 
  // esta función asegura que la IP es válida y que esté en el orden de bytes
  // correcto; guarda ese valor en &address->sin_addr.s_addr
  inet_pton(AF_INET, ip, &address->sin_addr.s_addr);
  
  return address;
}

// MAIN //
int runServer(int port, char *ip) {

  // cremos el file descriptor para el socket del servidor, es decir,
  // el identificador de la "antena" de red
  int serverSocketFD = createSocket();

  // creamos la dirección del servidor, para que otros puedan encontrarlo,
  // esta dirección contiene la ip y el puerto, lo manejamos como apuntador
  // para que se libere la memoria solita al terminar el programa
  struct sockaddr_in *serverAddress = createAddress(ip, port);

  // terminamoss el programa si no se pudo crear correctamente la dirección
  if (serverAddress == NULL)
    return 1;
  
  // establecemos la conexion entre el socket y nuestra dirección, de tal forma
  // que todo lo que llegue al puerto con el que se inicializó el servidor se
  // le pasa a este mismo; lo guardamos en un int para comprobar si hubo exito o no
  int succesfulConnection = bind(serverSocketFD,
				 (struct sockaddr *)serverAddress,
				 sizeof(*serverAddress));

  // comprobamos si la conexion tuvo exito o no
  if(succesfulConnection != 0){
    fprintf(stderr,"Error: Hubo un problema al iniciar el servidor.");
    free(serverAddress);
    return 1;
  } else
    printf("Servidor ejecutandose correctamente.");

  // Aquí hacemos que el socket del servidor deje de ser pasivo y sea un socket
  // de escucha, además determinamos cuantos clientes podemos dejar en fila de
  // espera, en este caso digamos que diez; como con bind guardamos el resultado
  // en un int, pues si es 0 tuvimos exito al poner el servidor en modo escucha,
  // si es -1 algo salió mal
  int listenResult = listen(serverSocketFD, 10);

  // Asignamos un espacio en la memoria donde guardaremos la información del cliente
  struct sockaddr_in clientAddress;

  // Calculamos el tamaño exacto de la estructura en memoria pues queremos reservar
  // ese espacio para que al llamar a la función accept() sepa donde escribir los
  // datos dados por el usuario sin corromper otros lados del sistema (i.e. sin que
  // haya problemas con la asignacion de memoria)
  int clientAddressSize = sizeof(struct sockaddr_in);

  // Creamos un socket para el cliente que se acaba de conectar, es decir, si
  // nadie se conecta el programa se detiene aquí esperando alguna conexion al
  // servidor y cuando eso suceda crea una linea privada entre el servidor y
  // el cliente conectado
  int clientSocketFD = accept(serverSocketFD,
			      (struct sockaddr *)&clientAddress,
			      &clientAddressSize);

  // Creamos un espacio en memoria para guardar el mensaje que escriba el cliente
  char message[512];

  // Aquí le estamos diciendo al sistema operativo que los datos que tenga
  // guardados del clientSocketFD (porque cada cliente tiene su socket) los copie
  // a la variable message y que solo copie 512 bytes incluso si el mensaje es
  // más largo
  recv(clientSocketFD, message, 512, 0);

  // Imprimimos el mensaje que envío el cliente
  printf("El cliente envió: %s\n", message);

  // Cerramos los sockets y liberamos memoria
  close(clientSocketFD);
  close(serverSocketFD);
  free(serverAddress);

  return 0;
  
}
