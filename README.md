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

# Notas
1. El middleware al iniciar y cada 30 segundos verifica el estado de los servidores enviando una solicitud http GET al endpoint `/info` y esperando un `200 OK`. De momento no hace validación del contenido obtenido, los servidores deben tener este endpoint y retornar un `200 OK` para ser considerados válidos y en línea.
2. Si un servidor no está registrado, pero intenta enviar una solicitud al middleware, al menos que sea `GET /info` o `PING` este recibirá un `401 Unauthorized`.
3. En el archivo de servidores tienen que estar registrados todos los servidores del cluster, tanto los que envían las solicitudes como los que lo reciben, el middleware puede detectar el servidor que envió la solicitud mediante el host de origen o la cabecera `X-Middleware-Sent-By`, los servidores no tienen que enviar esta cabecera a menos que sepan que la solicitud pasará por otros servidores que potencialmente reescriban el host de origen y este no esté registrado en el middleware. También pueden utilizar esta cabecera para saber qué servidor envió la actualización.
4. El middleware enviará las solicitudes con la cabecera `User-Agent` con el valor `middleware`, los servidores deben verificar que las peticiones no contengan este `User-Agent` para así poder reenviar la solicitud recibida al middleware

