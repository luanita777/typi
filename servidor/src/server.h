#ifndef SERVER_H
#define SERVER_H

#include <netinet/in.h>

int createSocket();
struct sockaddr_in* createAddress(char *ip, int port);
void runServer(int port, char *ip); 

#endif
