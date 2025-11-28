

# Reporte de cobertura de tests

## ¿Cómo están organizadas las pruebas?

- **Unitarias:**
	- Prueban la lógica de negocio de los servicios (usecase/task, usecase/tasklist) y los handlers HTTP de forma aislada.
	- Se encuentran en: `internal/usecase/task/service_test.go`, `internal/usecase/tasklist/service_test.go`, `internal/delivery/http/*_test.go`, y en `tests/unit/`.

- **Integración:**
	- Validan que la app funcione bien con la base de datos y el flujo completo de las rutas.
	- Están en: `tests/integration/`.

## ¿Qué se probó?

- Creación, consulta, actualización y borrado de tareas y listas.
- Validación de errores y respuestas HTTP.
- Lógica de negocio (por ejemplo, reglas de estados y prioridades).
- Integración real con la base de datos (insertar, consultar, borrar).

## ¿Cuánto cubre?

El último resultado de cobertura fue:

```
total: (statements) 63.2%
```

