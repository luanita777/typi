#ifndef SERVER_H
#define SERVER_H

#include <netinet/in.h>

int crearSocketServidor();
struct sockaddr_in* crearDireccionServidor(char *ip, int puerto);
void ejecutaServidor(int puerto, char *ip); 

#endif
