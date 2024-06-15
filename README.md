# Stori Technical Challenge

Este proyecto procesa un archivo de transacciones, las guarda en una base de datos SQLlite y envía un resumen por correo electrónico con el logo de la empresa.

## Configuración

1. Editar el archivo `txns.csv` con las transacciones.
2. Configurar el correo electrónico y las credenciales SMTP en tus variables de entorno.

## Compilación y Ejecución

### RUN Local

```sh
go run cmd/api/main.go
```

### TEST Local

```sh
go test ./...
```

### BUILD DOCKER Local

```sh
docker build -t lucasbellesi/stori-technical-challenge .
```

### RUN DOCKER Local

```sh
docker run --rm lucasbellesi/stori-technical-challenge
```