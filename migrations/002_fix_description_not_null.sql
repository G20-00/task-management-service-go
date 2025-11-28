-- Corrige la columna description para que nunca sea NULL y tenga valor por defecto
ALTER TABLE task_lists ALTER COLUMN description SET DEFAULT '';
UPDATE task_lists SET description = '' WHERE description IS NULL;
ALTER TABLE task_lists ALTER COLUMN description SET NOT NULL;
