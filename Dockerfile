# Usar una imagen base oficial de Go
FROM golang:1.22.4-alpine

# Instalar dependencias necesarias para cgo y SQLite
RUN apk add --no-cache gcc musl-dev

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el módulo Go y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar los archivos de la aplicación
COPY . .

# Establecer la variable de entorno para habilitar cgo
ENV CGO_ENABLED=1

# Compilar la aplicación
RUN go build -o /app/main ./cmd/api

# Comando para ejecutar el binario compilado
CMD ["/app/main"]