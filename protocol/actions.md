# Proceso

Una vez conectado con el middleware, los clientes comunicarán sus cambios con el mismo, y este propagará estos cambios a otros clientes. El protocolo para hacerlo es el siguiente:

```json
{
    "version": 2,
    "clientId": "<id pasado desde el servidor>",
    "action": "<put | del>",
    "entity": "<estudiante | profesor | asignatura | ciclo | matricula | rnota >",
    "args": {
        ...
    }
}
```

Dentro de `args` se encontrará la información modificada o eliminada, pudiendo ser `{ id_: "identificador global" }`, o la entidad modificada.

La respuesta de esta petición será:

```json
{
    "version": 2,
    "data": {
        "lastTimeStamp": 12340,
    }
}
```

# El valor de `lastTimeStamp`
Es importante que se almacene el valor de `lastTimeStamp`, ya que esto es una marca de tiempo de la última actualización efectuada, lo cual servirá para solicitad únicamente lo que ha cambiado a partir de ese entonces.

El valor de `lastTimeStamp` será siempre generado desde el middleware, no se espera que los clientes generen uno. Cuando un cliente recibe un mensaje solicitando una actualización o eliminación, el cliente debe esperar el campo adicional de `lastTimeStamp` y almacenarlo para su posterior uso. 

El `lastTimeStamp` retornado en las respuestas del cliente se ignorarán, así que no es necesario que estos retornen dicha información.


# Ejemplos
## Modificar Profesor

Solicitud:
```json
{
    "version": 2,
    "clientId": "foo123",
    "action": "put",
    "entity": "profesor",
    "args": {
        "nombre": "Mario"
    }
}
```

Respuesta:
```json
{
    "version": 2,
    "data": {
        "lastTimeStamp": 12340,
    }
}
```

Solicitud reenviada desde el middleware:
```json
{
    "version": 2,
    "clientId": "foo123",
    "action": "put",
    "entity": "profesor",
    "args": {
        "lastTimeStamp": 12340,
        "nombre": "Mario"
    }
}
```

Respuesta:
```json
{
    "version": 2,
    "data": { }
}
```

## Eliminar un estudiante

Solicitud:
```json
{
    "version": 2,
    "clientId": "foo123",
    "action": "del",
    "entity": "estudiante",
    "args": {
        "id": "estu-123"
    }
}
```

Respuesta:
```json
{
    "version": 2,
    "data": {
        "lastTimeStamp": 12340,
    }
}
```

Solicitud reenviada desde el middleware:
```json
{
    "version": 2,
    "clientId": "foo123",
    "action": "del",
    "entity": "estudiante",
    "args": {
        "lastTimeStamp": 12340,
        "id": "estu-123"
    }
}
```

Respuesta:
```json
{
    "version": 2,
    "data": { }
}
```
