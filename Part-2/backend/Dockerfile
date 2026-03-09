# Multi-stage build para optimizar el tamaño de la imagen
FROM golang:1.21-alpine AS builder

# Instalar dependencias de compilación
RUN apk add --no-cache git

# Establecer directorio de trabajo
WORKDIR /app

# Copiar go mod files
COPY go.mod go.sum* ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Etapa final - imagen mínima
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el binario compilado
COPY --from=builder /app/main .

# Exponer puerto
EXPOSE 8080

# Ejecutar la aplicación
CMD ["./main"]
