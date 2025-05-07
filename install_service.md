

```cmd
go build pkg/server -o "<path_to_the_service_executable>" 

# Como administrador
sc.exe create Middleware binPath= "<path_to_the_service_executable>" 
```
