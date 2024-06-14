# Stori Technical Challenge

Este proyecto procesa un archivo de transacciones y envía un resumen por correo electrónico.

## Configuración

1. Editar el archivo `txns.csv` con las transacciones.
2. Configurar el correo electrónico y las credenciales SMTP en el archivo `.env`.

## Compilación y Ejecución

### Local

```sh
docker build -t stori-technical-challenge .

docker run --rm stori-technical-challenge