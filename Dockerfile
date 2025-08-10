# Build stage
FROM golang:1.24-alpine AS builder

# Instalar dependencias del sistema
RUN apk add --no-cache git ca-certificates tzdata

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Construir la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Production stage
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Crear usuario no-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Establecer directorio de trabajo
WORKDIR /app

# Copiar el binario desde el stage de build
COPY --from=builder /app/main .

# Copiar archivos de configuración
COPY --from=builder /app/env.example ./env.example

# Cambiar propiedad de archivos
RUN chown -R appuser:appgroup /app

# Cambiar al usuario no-root
USER appuser

# Exponer puerto
EXPOSE 8080

# Variables de entorno por defecto
ENV ENVIRONMENT=production
ENV SERVER_PORT=8080
ENV SERVER_HOST=0.0.0.0

ENV DB_HOST=find-my-friend-db.database.windows.net
ENV DB_PORT=1433
ENV DB_NAME=find-my-friend
ENV DB_USER=db_admin
ENV DB_PASSWORD=Isabella1506*
ENV B_SSL_MODE=disable

ENV JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
ENV JWT_EXPIRATION_HOURS=24

ENV CLOUDINARY_CLOUD_NAME=dw9e57leg
ENV CLOUDINARY_API_KEY=253523817534744
ENV CLOUDINARY_API_SECRET=2K5NwA4y4vKtUdKaXKV1uZjtCr0

# Comando para ejecutar la aplicación
CMD ["./main"] 