## Requisito para Makefile en Windows

Para usar los comandos automáticos del Makefile en Windows, primero instala make con Chocolatey:

```powershell
choco install make
```

Luego, cierra y vuelve a abrir la terminal para que el comando `make` esté disponible.
## Ejecutar todas las tareas del Makefile y ver cobertura

Puedes ejecutar todas las tareas principales (formato, lint, tests y cobertura) en una sola línea:

```powershell
make format && make lint && make test && go test ./... -coverprofile=coverage && go tool cover -func=coverage
```

Esto dejará el código formateado, limpio, probado y mostrará el resumen de cobertura.
## Uso de Makefile

Puedes automatizar tareas comunes usando el Makefile incluido. Los comandos principales son:

- Formatear código:
    ```bash
    make format
    ```
- Linting:
    ```bash
    make lint
    ```
- Ejecutar tests:
    ```bash
    make test
    ```

Si quieres ejecutar todas las tareas (formato, lint, tests y ver cobertura) en una sola línea, usa:

```powershell
make format && make lint && make test && go test ./... -coverprofile=coverage && go tool cover -func=coverage
```

Esto dejará el código formateado, limpio, probado y mostrará el resumen de cobertura.

Esto facilita mantener el código limpio y ejecutar pruebas rápidamente.
## Formateo de código

Para asegurar el formato correcto del código, ejecuta:

```bash
gofmt -w .
goimports -w .
```

Esto aplicará el formato estándar de Go y organizará automáticamente los imports en todos los archivos del proyecto.
## Limpiar archivo de cobertura

Si deseas eliminar el archivo de cobertura generado, ejecuta:

```powershell
del coverage
```

Esto borra el archivo `coverage` generado por los tests.
## Linter

Para asegurar la calidad y el estilo del código, ejecuta el siguiente comando:

```bash
golangci-lint run
```

Esto revisará el proyecto con las reglas configuradas y mostrará advertencias o errores de estilo, formato y buenas prácticas.
# Task Management Service

API REST para gestión de tareas construida con Go y Fiber.

## **_Construido con_** 

<div style="text-align: left">
    <p>
        <a href="https://go.dev/" target="_blank"> 
            <img alt="Go" src="https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Blue.png" height="60" width="auto" style="vertical-align: bottom;">
        </a>
        <a href="https://docs.docker.com/" target="_blank"> 
            <img alt="Docker" src="https://www.docker.com/wp-content/uploads/2022/03/vertical-logo-monochromatic.png" height="60" width="auto" style="vertical-align: bottom;">
        </a>
        <a href="https://www.postgresql.org/" target="_blank">
            <img src="https://wiki.postgresql.org/images/a/a4/PostgreSQL_logo.3colors.svg" height="60" width="auto" alt="PostgreSQL" style="vertical-align: bottom;">
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
### Base de datos y migraciones

```bash
# Levantar postgres (opcional si no usas docker-compose)
docker run --name task-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=task_management -p 5432:5432 -d postgres:15

# Crear la BD (si no existe)
docker exec -it task-postgres psql -U postgres -c "CREATE DATABASE task_management;"

# Aplicar migraciones (requiere migrate instalado)

Este comando aplica las migraciones de base de datos usando la herramienta [migrate](https://github.com/golang-migrate/migrate). Es necesario para crear o actualizar las tablas según los archivos SQL en la carpeta `migrations`.

**Descarga para Windows:**
- Ve a: https://github.com/golang-migrate/migrate/releases
- Busca la sección "Assets" de la última versión.
- Descarga el archivo `migrate.windows-amd64.zip`.
- Extrae el archivo `migrate.exe` y colócalo en la raíz del proyecto o en una carpeta incluida en tu PATH.

**Ejecuta:**
```powershell
./migrate.exe -database "postgres://postgres:postgres@127.0.0.1:5432/task_management?sslmode=disable" -path ./migrations up
```
Esto aplicará todas las migraciones pendientes a la base de datos.
```
### Docker

```bash
# Levantar toda la app y la base de datos con Docker Compose
docker-compose up --build
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
```bash
# Correr localmente
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

## Pruebas y cobertura

Asegúrate de tener todas las dependencias instaladas y la base de datos corriendo. Para ejecutar las pruebas y ver el reporte de cobertura, usa los siguientes comandos en PowerShell:

```powershell
go test ./... -coverprofile=coverage
go tool cover -func=coverage
```

Esto generará el archivo `coverage` y mostrará el resumen de cobertura en consola.

## **_Autor_** 

<div style="text-align: left">
    <a href="https://github.com/G20-00" target="_blank"> <img alt="G20-00" src="https://images.weserv.nl/?url=https://avatars.githubusercontent.com/u/70019070?v=4&h=60&w=60&fit=cover&mask=circle"></a>
</div>
