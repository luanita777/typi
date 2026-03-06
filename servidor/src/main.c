#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <limits.h>
#include <arpa/inet.h>
#include "server.h"

int main(int numArgs, char *args[]){
  if(numArgs != 3){
    fprintf(stderr, "Uso: %s <ip> <puerto>\n", args[0]);
    return 1;
  }

  //Validamos IP
  char *received_ip = args[1];

  struct in_addr addr_aux;
  if (inet_pton(AF_INET, received_ip, &addr_aux) != 1) {
    fprintf(stderr, "Error: '%s' no es una dirección IPv4 válida.\n", received_ip);
    return 1;
  }

  //Validamos Puerto
  char *endptr;
  errno = 0;
  long puertoL = strtol(args[2], &endptr, 10);

  if(args[2] == endptr){
    fprintf(stderr, "Error: el puerto no es un número válido.\n");
    return 1;
  }

  if (*endptr != '\0') {
    fprintf(stderr, "Error: el puerto contiene caracteres inválidos.\n");
    return 1;
  }

  if ((errno == ERANGE && (puertoL == LONG_MAX || puertoL == LONG_MIN))) {
    fprintf(stderr, "Error: el número está fuera de rango.\n");
    return 1;
  }
  
  int puerto = (int) puertoL;
  
  printf("[INFO] IP y puerto validados. Iniciando...\n");
  runServer(puerto, received_ip);
  return 0;
}
