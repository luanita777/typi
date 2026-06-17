# Typi – Chat tipo IRC

**Typi** es un proyecto de chat que permite que múltiples usuarios se conecten a un servidor central y se comuniquen en tiempo real. 

Este repositorio contiene al **servidor**, además de un `docker-compose.yml` para levantar servicios de forma sencilla.

---

## Requisitos para usarlo

- [Docker](https://www.docker.com/get-started)  
- [Docker Compose](https://docs.docker.com/compose/install/)  
- Go 1.23+ para el servidor 

---
## Levantar el servidor con Docker
1. Elegir un puerto libre en el que queramos que escuche el servidor, por ejemplo `1234`.
2. Desde la raíz del proyecto, ejecutar:
      `sudo PUERTO=1234 docker compose up --build

