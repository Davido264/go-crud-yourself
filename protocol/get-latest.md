# Proceso

El cliente o servidor solicita la información a un determinado tiempo mediante el siguiente mensaje:

```json
{
    "version": 2,
    "clientId": "<id pasado desde el servidor>",
    "action": "get",
    "entity": "<estudiante | profesor | asignatura | ciclo | matricula | rnota >",
    "args": {
        "lastTimeStamp": 123
    }
}
````

Respuesta:

```json
{
    "version": 2,
    "data": {
        "lastTimeStamp": 12340,
        "estudiante": [...]
    }
}
```

# Ejemplo

El cliente solicita los datos más actualizados de los **estudiantes**.

```json
{
    "version": 2,
    "clientId": "foo123",
    "action": "get",
    "entity": "estudiante",
    "args": {
        "lastTimeStamp": 1746047938
    }
}
```

El valor de `lastTimeStamp` es enviado con cada actualización, lo que significa que este valor corresponde a la última actualización que se realizó. Si no se ha realizado ninguna, se puede pasar `null`

```json
{
    "version": 2,
    "clientId": "foo123",
    "action": "get",
    "entity": "estudiante",
    "args": {
        "lastTimeStamp": null 
    }
}
```

La respuesta será la información solicitada:

```json
{
    "version": 2,
    "data": {
        "lastTimeStamp": 1746048118,
        "estudiante": [
            {
                "nombre": "Mario"
            }
        ]
    }
}
```

Esta puede ser un array vació en caso de que ya se encuentre actualizado:

```json
{
    "version": 2,
    "data": {
        "lastTimeStamp": 1746048118,
        "estudiante": []
    }
}
```

# Contenido del objeto `data`
La información en un objeto conteniendo 2 atributos: `lastTimeStamp` y la entidad solicitada, por ejemplo:
- Si el cliente envía "profesor" en `entity`, recibirá `{ "lastTimeStamp": 1746048118, "profesor": [...] }`
- Si el cliente envía "estudiante" en `entity`, recibirá `{ "lastTimeStamp": 1746048118, "estudiante": [...] }`
- Si el cliente envía "matricula" en `entity`, recibirá `{ "lastTimeStamp": 1746048118, "matricula": [...] }`
