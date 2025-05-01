# Proceso

El cliente o servidor envía un mensaje para obtener el estado

```json
{
    "version": 2,
    "clientId": "<id pasado desde el servidor>",
    "action": "get",
    "entity": "status",
    "args": { }
}
````

Respuesta:

```json
{
    "version": 2,
    "data": {
        "lastTimeStamp": 12340
    }
}
```

# Ejemplo

El servidor solicita saber cuál es el estado actual del cliente.

```json
{
    "version": 2,
    "clientId": "foo123",
    "action": "get",
    "entity": "status",
    "args": {
        "lastTimeStamp": 1746047938
    }
}
```

La respuesta será la información solicitada:

```json
{
    "version": 2,
    "data": {
        "lastTimeStamp": 1746048118,
    }
}
```
