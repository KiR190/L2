-- Удаляем индексы
DROP INDEX IF EXISTS idx_events_user_date;
DROP INDEX IF EXISTS idx_events_date;
DROP INDEX IF EXISTS idx_events_user_id;

-- Удаляем триггер
DROP TRIGGER IF EXISTS set_updated_at ON events;

-- Удаляем функцию
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Удаляем таблицу
DROP TABLE IF EXISTS events;
