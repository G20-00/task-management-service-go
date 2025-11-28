
# DECISION_LOG.md

## ¿Por qué Fiber?
Me fui por Fiber porque es rápido, fácil de usar y tiene buena documentación. Además, su rendimiento es excelente para APIs REST.

## ¿Y el ORM?
No usé un ORM pesado. Mejor el driver de Postgres (`lib/pq`) con `database/sql`, así tengo control total sobre las queries y no hay magia rara.

## Cosas técnicas que decidí
- Separé el código en capas (domain, usecase, delivery, infra) para que no sea un espagueti.
- Uso migraciones con golang-migrate, así la base de datos no se descontrola.
- Todo corre en Docker, para que sea fácil de levantar en cualquier lado.
- Configuración con variables de entorno, para no hardcodear nada.
- Hay tests unitarios y de integración para asegurar que todo funcione bien.
- Linter y formateo automático con golangci-lint para mantener el código limpio.

## Cosas que me faltan o podría mejorar
- Subir la cobertura de tests en algunos archivos.
- Meter context.Context en más lados.
- Manejo de errores más detallado (códigos HTTP específicos).
- Documentar los endpoints con Swagger.
- CI/CD listo con GitHub Actions: cada push corre los tests y lint automáticamente.
