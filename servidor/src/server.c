#include <stdio.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <pthread.h>

struct SocketAceptado{
  int socketClienteFD;
  struct sockaddr_in direccion;
  int aceptadoCorrectamente;
  int error;
};

/* ========== FUNCIONES AUXILIARES ========== */

 /**
  * @brief Crea el socket principal del servidor.
  *
  * @details
  * Inicializa un socket utilizando IPv4 (AF_INET) y comunicación orientada
  * a conexión mediante TCP (SOCK_STREAM). Este socket funcionará como el
  * punto de entrada para todas las conexiones de clientes al servidor.
  *
  * @return int
  * File descriptor del socket creado.
  *
 */
  int crearSocketServidor(){
    return socket(AF_INET, SOCK_STREAM, 0);
  }

  /**
   * @brief Construye la dirección del servidor.
   *
   * @details
   * Reserva memoria para una estructura sockaddr_in, la limpia asignandola
   * toda a cero con memset() y la inicializa con la dirección IP y el 
   * el puerto especificados (usamos -> porque address es un apuntador)
   * La IP se convierte al formato binario requerido por el sistema
   * operativo mediante inet_pton.
   *
   * Esta estructura será utilizada posteriormente por bind() para asociar
   * el socket del servidor a una dirección de red específica.
   *
   * @param ip Dirección IP en formato string (ej: "127.0.0.1")
   * @param puerto Puerto en el que el servidor escuchará conexiones
   *
   * @return struct sockaddr_in*
   * Apuntador a la estructura con la dirección configurada.
   *
  */
   struct sockaddr_in* crearDireccionServidor(char *ip, int puerto){
     struct sockaddr_in *direccion = malloc(sizeof(struct sockaddr_in));
     memset(direccion, 0, sizeof(struct sockaddr_in));

     direccion->sin_family = AF_INET;
     direccion->sin_port = htons(puerto);

     inet_pton(AF_INET, ip, &direccion->sin_addr.s_addr);
     
     return direccion;
   }


   /**
    * @brief Acepta una nueva conexión entrante de un cliente.
    *
    * @details
    * Asignamos un espacio en la memoria donde guardaremos la información
    * del cliente, calculamos el tamaño exacto de la estructura en memoria
    * pues queremos reservar ese espacio para que al llamar a la función
    * accept() no haya problemas con la asignacion de memoria.
    *
    * Luego, utilizamos la función accept() para bloquear la ejecución hasta que 
    * un cliente intente conectarse al servidor. Cuando esto ocurre, el sistema
    * operativo toma una conexión pendiente de la cola del socket de escucha y
    * crea un nuevo socket dedicado exclusivamente a la comunicación con ese
    * cliente. Este nuevo socket es independiente del socket del servidor,
    * permitiendo que el servidor continúe aceptando más conexiones mientras
    * cada cliente se comunica a través de su propio descriptor.
    *
    * Finalmente, la información relevante de la conexión (dirección del cliente,
    * descriptor del socket y estado de la conexión) se encapsula dentro de una
    * estructura SocketAceptado para facilitar su manejo dentro del programa.
    *
    * @param socketServidorFD File descriptor del socket del servidor
    *
    * @return struct SocketAceptado*
    * Estructura con la información del cliente conectado.
    *
   */
   struct SocketAceptado* aceptarConexionCliente(int socketServidorFD){
     struct sockaddr_in direccionCliente;
     int tamañoDireccionCliente = sizeof(struct sockaddr_in);

     int socketClienteFD = accept(socketServidorFD,
			      (struct sockaddr *)&direccionCliente,
			       &tamañoDireccionCliente);

     struct SocketAceptado* socketAceptado = malloc(sizeof(struct SocketAceptado));
     socketAceptado->direccion = direccionCliente;
     socketAceptado->socketClienteFD = socketClienteFD;
     socketAceptado->aceptadoCorrectamente = socketClienteFD > 0;
     
     if(!socketAceptado->aceptadoCorrectamente)
       socketAceptado->error = socketClienteFD;

     return socketAceptado;
   }

  /**
   * @brief Función ejecutada por cada hilo para gestionar la comunicación
   * con un cliente conectado.
   *
   * @details
   * La firma `void* funcion(void*)` es requerida por la biblioteca pthread.
   * Esto se debe a que C no tiene genéricos, por lo que pthread utiliza
   * `void*` para poder recibir cualquier tipo de dato como argumento.
   *
   * En este caso, el hilo recibe un apuntador a un entero que contiene el
   * descriptor del socket del cliente. Como pthread lo recibe como `void*`,
   * es necesario convertirlo nuevamente a `int*` haciendole una audición 
   * y dereferenciarlo para recuperar el descriptor real del socket.
   *
   * Una vez obtenido el socket del cliente, el hilo entra en un ciclo donde
   * recibe mensajes enviados por el cliente usando recv() y los imprime
   * en el servidor. Si el cliente cierra la conexión o ocurre un error,
   * el ciclo termina y el socket se cierra.
   *
   * @param arg Apuntador genérico (void*) que en realidad contiene un int*
   *            con el descriptor del socket del cliente.
   *
   * @return void* No regresa ningún valor significativo.
   */
   void* gestionarCliente(void *arg){
     int socketFD = *(int*)arg;
     free(arg);
     char mensaje[512];
     
     while(1){
       
       int mensajeRecibidoExitosamente = recv(socketFD, mensaje, 512, 0);
       
       if(mensajeRecibidoExitosamente > 0){
	 mensaje[mensajeRecibidoExitosamente] = '\0';
	 printf("El cliente envió: %s\n", mensaje);
       } else {
	 if(mensajeRecibidoExitosamente == -1)
	   fprintf(stderr, "Error: Sucedió un problema al recibir los datos.\n");
	 else
	   printf("\nEl cliente cerró la conexión.\n");
	 break;
       } 
     }
     
     close(socketFD);     
   }

 /**
  * @brief Crea un hilo para gestionar un cliente.
  *
  * @details
  * Esta función inicia un nuevo hilo usando la biblioteca pthreads.
  * El hilo ejecutará la función gestionarCliente(), permitiendo que
  * cada cliente sea atendido de manera concurrente sin bloquear
  * el servidor principal.
  *
  * @param socketCliente Estructura con la información del cliente aceptado.
  *
 */
  void crearHiloCliente(struct SocketAceptado *socketCliente){
    pthread_t idHilo;
    int *socketFD = malloc(sizeof(int));
    *socketFD = socketCliente->socketClienteFD;
    pthread_create(&idHilo, NULL, gestionarCliente, socketFD);
    pthread_detach(idHilo);
  }

/**
  * @brief Mantiene al servidor escuchando nuevas conexiones.
  *
  * @details
  * Ejecuta un ciclo infinito donde continuamente se aceptan nuevas
  * conexiones mediante aceptarConexionCliente(). Cada cliente aceptado
  * se maneja en un hilo independiente.
  *
  * Esto permite que el servidor maneja múltiples clientes de forma
  * concurrente.
  *
  * @param socketServidorFD File descriptor del socket del servidor.
  *
 */
  void escucharConexiones(int socketServidorFD){
    while (1) {
      
      struct SocketAceptado* cliente = aceptarConexionCliente(socketServidorFD);
      
      if (cliente->aceptadoCorrectamente) 
	crearHiloCliente(cliente);
      else 
	fprintf(stderr, "Error aceptando conexión\n");
      
    }
  }
  

/* =========== MAIN =========== */

 /**
  * @brief Inicializa y ejecuta el servidor.
  *
  * @details
  * Esta función se encarga de coordinar todo el proceso de arranque del servidor:
  *
  * 1. Crea el socket del servidor.
  * 2. Construye la dirección IP y puerto donde escuchará.
  * 3. Asocia el socket a la dirección usando bind().
  * 4. Activa el modo de escucha mediante listen().
  * 5. Comienza a aceptar conexiones entrantes.
  *
  * Si ocurre un error durante la inicialización, el servidor termina
  * devolviendo un código de error.
  *
  * @param puerto Puerto donde el servidor escuchará conexiones
  * @param ip Dirección IP donde el servidor será accesible
  *
  * @return int
  * 0 si el servidor finaliza correctamente, 1 si ocurre un error.
  *
 */
  int ejecutaServidor(int puerto, char *ip) {
    
     int socketServidorFD = crearSocketServidor();
     struct sockaddr_in *direccionServidor = crearDireccionServidor(ip, puerto);   
     if (direccionServidor == NULL)
       return 1;
     
     int conexionExitosa = bind(socketServidorFD,
				(struct sockaddr *)direccionServidor,
				sizeof(*direccionServidor));
     
     if(conexionExitosa != 0){
       fprintf(stderr,"Error: Hubo un problema al iniciar el servidor.");
       free(direccionServidor);
       return 1;
     } else
       printf("Servidor ejecutandose correctamente...\n\n");
     
     listen(socketServidorFD, 10);
     
     escucharConexiones(socketServidorFD);
     
     close(socketServidorFD);
     free(direccionServidor);
     
     return 0;    
}
