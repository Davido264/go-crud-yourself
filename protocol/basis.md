# Introducción

La comunicación entre los backends con el middleware se realiza mediante json enviados por websockets. Los mensajes que puede recibir y enviar un cliente son 2, uno de éxito y uno de error.

Mensaje de éxito:

```json
{
    "version": 2,
    "data": {
        ...
    }
}
```

Mensaje de error:

```json
{
    "version": 2,
    "errno": "<mensaje de error>"
}
```

# Comunicación bidireccional
Tanto el cliente como el servidor pueden enviar y recibir mensajes, así que se debe implementar tanto la interpretación de mensajes como la encodificación de los mensajes al servidor.
Un cliente debe esperar la recepción los mensajes definidos en los todos estos archivos.

# Conectarse al servidor

Cuando un cliente se conecta por primera vez al servidor, recibirá su `clientId` de la siguiente forma:

```json
{
    "version": 2,
    "data": {
        "clientId": "client-1234"
    }
}
```

Este `clientId` tiene que ser enviado en cada mensaje para identificar al cliente. El middleware por su parte, cuando envíe mensajes a los clientes, utilizará el `clientId` del cliente. Esto debido a que no hay muchos middleware para un cliente, pero sí muchos clientes para un middleware, lo que significa que a diferencia del middleware, este campo no tiene utilidad en mensajes enviados desde el middleware hasta el cliente. Se espera que los clientes ignoren este campo.

# Entidades que se pueden solicitar
- `estudiante`: Estudiantes
- `profesor`: Profesores
- `asignatura`: Asignaturas
- `ciclo`: Ciclos
- `matricula`: Matrículas
- `rnota`: Registros de Notas
- `status`: Estado actual del nodo

# Acciones que se pueden realizar
- `get`: Obtener el valor o valores de la entidad especificada
- `put`: Crear o modificar la entidad especificada
- `del`: Eliminar la entidad especificada

# Mensajes de error
Los mensajes de error que se pueden recibir son:
- `ERRNO_INVALID_FIELD`: Indica que el mensaje que se envió no tiene los atributos correctamente definidos
- `ERRNO_INVALID_FORMAT`: Indica que el mensaje que se envió está en un formato incorrecto
- `ERRNO_NOT_ALLOWED`: Indica que la acción no está disponible para la  entidad. Este error se obtendrá principalmente de realizar una acción que no sea `get` a `status`
- `ERRNO_INVALID_ARGS`: Indica que los parámetros enviados no son los esperados para determinada acción
- `ERRNO_INVALID_PROTOCOL_VERSION`: Indica que la versión del protocolo es incompatible

Se incluirán más a medida que avance el desarrollo, así que se puede esperar que se agreguen más códigos de error
