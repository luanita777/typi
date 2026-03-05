#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <limits.h>

int main(int numArgs, char *args[]){
  if(numArgs < 2 || numArgs > 2){
    fprintf(stderr, "Uso: %s <puerto>\n", args[0]);
    return 1;
  }

  char *endptr;
  errno = 0;
  long puertoL = strtol(args[1], &endptr, 10);

  if(args[1] == endptr){
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
  
  printf("[INFO] Servidor iniciando en puerto %d\n", puerto);
  
  return 0;
}
