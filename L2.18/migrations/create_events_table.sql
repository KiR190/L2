-- Создание таблицы events
CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    event_date DATE NOT NULL,
    title TEXT NOT NULL CHECK (char_length(title) > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON events
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_events_user_date ON events (user_id, event_date);
CREATE INDEX IF NOT EXISTS idx_events_date ON events (event_date);
CREATE INDEX IF NOT EXISTS idx_events_user_id ON events (user_id);

INSERT INTO events (user_id, event_date, title)
VALUES
    (1, '2025-10-19', 'Посещение врача'),
    (1, '2025-10-21', 'Встреча с командой'),
    (2, '2025-10-22', 'Поход в горы'),
    (1, '2025-11-01', 'День рождения друга');