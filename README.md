# Task Management Service

API REST para gestión de tareas construida con Go y Fiber.

## **_Construido con_** 

<div style="text-align: left">
    <p>
        <a href="https://go.dev/" target="_blank"> 
            <img alt="Go" src="https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Blue.png" height="80" width="80" style="vertical-align: bottom;">
        </a>
        <a href="https://docs.docker.com/" target="_blank"> 
            <img alt="Docker" src="https://miro.medium.com/v2/resize:fit:453/1*QVFjsW8gyIXeCUJucmK4XA.png" height="80" width="80" style="vertical-align: bottom;">
        </a>
        <a href="https://www.postgresql.org/" target="_blank">
            <img src="https://www.postgresql.org/media/img/about/press/elephant.png" height="80" width="80" alt="PostgreSQL" style="vertical-align: bottom;">
        </a>
    </p>
</div>

## Setup

### Base de datos

```bash
# Levantar postgres
docker run --name task-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=task_management -p 5432:5432 -d postgres:15

# Crear la BD
docker exec -it task-postgres psql -U postgres -c "CREATE DATABASE task_management;"

# Crear tablas
Get-Content database_schema.sql | docker exec -i task-postgres psql -U postgres -d task_management
```

### Variables de entorno

Crear `.env`:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=task_management
DB_SSLMODE=disable
```

### Correr

```bash
go run cmd/api/main.go
```

## Endpoints

**TaskLists**
- POST `/api/lists` - Crear lista
- GET `/api/lists` - Ver todas
- GET `/api/lists/:id` - Ver una
- PUT `/api/lists/:id` - Actualizar
- DELETE `/api/lists/:id` - Eliminar

**Tasks**
- POST `/api/tasks` - Crear tarea
- GET `/api/tasks` - Ver todas
- GET `/api/tasks/:id` - Ver una
- PUT `/api/tasks/:id` - Actualizar
- DELETE `/api/tasks/:id` - Eliminar

## Ejemplos

```powershell
# Crear lista
Invoke-RestMethod -Uri "http://localhost:8080/api/lists" -Method Post -Body '{"name":"Lista 1","description":"test"}' -ContentType "application/json"

# Crear tarea
Invoke-RestMethod -Uri "http://localhost:8080/api/tasks" -Method Post -Body '{"list_id":"ID","title":"Tarea 1","priority":"high"}' -ContentType "application/json"
```

Estados: `pending`, `in-progress`, `completed`  
Prioridades: `low`, `medium`, `high`

## **_Autor_** 

<div style="text-align: left">
    <a href="https://github.com/G20-00" target="_blank"> <img alt="G20-00" src="https://images.weserv.nl/?url=https://avatars.githubusercontent.com/u/70019070?v=4&h=60&w=60&fit=cover&mask=circle"></a>
</div>
