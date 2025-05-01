# Cómo correrlo

1. Crear un archivo, este puede ser de cualquier nombre. Su contenido será la lista de urls de servidores que se registrarán en el middleware
    ```
    localhost:8080 misitio.com https://holamundo.com
    localhost:8081 sitio2.com
    ```
    En ese ejemplo se registrarán 2 servidores, el primero con las direcciones ip o urls `localhost:8080`, `misitio.com`, y `https://holamundo.com`, y el segundo con las direcciones ip o urls `localhost:8081`, `sitio2.com`

2. Ejecutar el siguiente comando en el directorio del proyecto
    ```sh
    go run . -config <archivo de servidores> -port <puerto>
    ```
    En donde `<archivo de servidores>` es el archivo creado con la lista de servidores, y `<puerto>` es el puerto en el que el middleware escuchará peticiones

# Cómo comunicarse con el?
El protocolo de comunicación está basado en websockest y json, para saber sobre este, se deben ver los siguientes documentos en orden:
1. [Introducción](./protocol/basis.md)
2. [Obtener el Último estado](./protocol/get-latest.md)
3. [Notificar cambios](./protocol/actions.md)
