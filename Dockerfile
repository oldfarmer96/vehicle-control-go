# ==========================================
# Etapa 1: Builder (Compilación)
# ==========================================
FROM golang:1.26.1-alpine AS builder

# Instalar dependencias del sistema necesarias
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Crear un usuario no-root por seguridad
ENV USER=appuser
ENV UID=10001
RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid "${UID}" \
  "${USER}"

WORKDIR /app

# Aprovechar el caché de capas de Docker descargando primero los módulos
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copiar el código fuente
COPY . .

# Compilar el binario (Altamente optimizado)
# CGO_ENABLED=0 -> Produce un binario estático sin dependencias de C
# -ldflags="-w -s" -> Reduce drásticamente el tamaño eliminando info de debug
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o api-server ./cmd/api/main.go

# ==========================================
# Etapa 2: Imagen Final (Producción)
# ==========================================
# 'scratch' es una imagen completamente vacía. Máxima seguridad.
FROM scratch

# Importar certificados de seguridad (para HTTPS) y zonas horarias
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Importar el usuario no-root
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Traer solo el binario compilado de la etapa anterior
COPY --from=builder /app/api-server /api-server

# Forzar la ejecución de la app con el usuario sin privilegios
USER appuser:appuser

# Variables por defecto (se sobrescriben con docker-compose)
ENV PORT=8080
ENV APP_ENV=production

EXPOSE 8080

# Ejecutar el binario
ENTRYPOINT ["/api-server"]
