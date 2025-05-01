# V2
## Cómo correrlo

1. Crear un archivo `/etc/config.json` o `./config.json` e incluir las configuraciones siguiendo el siguiente esquema:
    ```json
    {
        "port": 8069,
        "protocolVersion": 2,
        "chsize": 512,
        "servers": [
            {
                "alias": "<Alias del Servidor para identificarlo de mejor forma>",
                "address": [
                    "<dirección ip del servidor>",
                    ...
                ]
            }
        ]
    }
    ```

2. Ejecutar el siguiente comando en el directorio del proyecto
    ```sh
    go run ./pkg/server
    ```

## Cómo comunicarse con el?
El protocolo de comunicación está basado en websockest y json, para saber sobre este, se deben ver los siguientes documentos en orden:
1. [Introducción](./protocol/basis.md)
2. [Obtener el Último estado](./protocol/get-latest.md)
2. [Obtener estado del nodo](./protocol/get-status.md)
3. [Notificar cambios](./protocol/actions.md)
