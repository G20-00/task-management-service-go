-- Asegura que no haya NULLs en description en task_lists
UPDATE task_lists SET description = '' WHERE description IS NULL;
