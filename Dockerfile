# Usar una imagen base oficial de Go
FROM golang:1.18-alpine

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el módulo Go y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar los archivos de la aplicación
COPY . .

# Compilar la aplicación
RUN go build -o /app/main ./cmd/myproject

# Comando para ejecutar el binario compilado
CMD ["/app/main"]
