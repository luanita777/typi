# =========================
# Configuración
# =========================
CC = gcc
CFLAGS = -Wall -Wextra -Werror -std=c11
SRC_DIR = src
SRCS = $(wildcard $(SRC_DIR)/*.c)
OUT = server

# =========================
# Regla principal
# =========================
all: $(OUT)

$(OUT): $(SRCS)
	$(CC) $(CFLAGS) $(SRCS) -o $(OUT)

# =========================
# Limpieza
# =========================
clean:
	rm -f $(OUT)

# =========================
# Ejecutar
# =========================
run: $(OUT)
	./$(OUT) 8080

