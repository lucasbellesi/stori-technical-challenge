# Stori Technical Challenge

Este proyecto procesa un archivo de transacciones y envía un resumen por correo electrónico.

## Configuración

1. Crear un archivo `transactions.csv` con las transacciones.
2. Configurar el correo electrónico y las credenciales SMTP en `pkg/email/email.go`.

## Compilación y Ejecución

### Local

```sh
go run cmd/api/main.go