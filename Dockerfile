# ===================================
# STAGE 1 — Build Servidor (C)
# ===================================
FROM ubuntu:22.04 AS server-build

RUN apt update
RUN apt install -y build-essential
RUN apt install -y make
RUN apt install -y valgrind
RUN rm -rf /var/lib/apt/lists/*

WORKDIR /build/server

COPY servidor/ .

RUN make
# RUN make test   # cuando tengas tests


# ===================================
# STAGE 2 — Build Cliente (Rust)
# ===================================
FROM ubuntu:22.04 AS client-build

RUN apt update
RUN apt install -y curl
RUN apt install -y build-essential
RUN rm -rf /var/lib/apt/lists/*

RUN curl https://sh.rustup.rs -sSf | sh -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"

WORKDIR /build/client

COPY cliente/ .

RUN cargo test
RUN cargo build --release


# ===================================
# STAGE 3 — Runtime (Producción)
# ===================================
FROM ubuntu:22.04 AS runtime

WORKDIR /app

COPY --from=server-build /build/server/server ./server
COPY --from=client-build /build/client/target/release/cliente ./client

CMD ["./server", "8080"]