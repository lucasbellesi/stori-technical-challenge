# Stori Technical Challenge

Este proyecto procesa un archivo de transacciones y envía un resumen por correo electrónico.

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