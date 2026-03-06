#ifndef SERVER_H
#define SERVER_H

#include <netinet/in.h>

int crearSocket();
struct sockaddr_in* crearDireccion(char *ip, int puerto);
void ejecutaServidor(int puerto, char *ip); 

#endif
