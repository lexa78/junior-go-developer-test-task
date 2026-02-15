CREATE
EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE subscriptions
(
--    делаю uuid, а не int autoincrement, т.к. у id user-а uuid,
--    чтобы все сущности в БД были с ключем одинакового типа
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--    хочу козырнуть своими знаниями о Нормальных формах ))
--    при создании сервиса с нуля, я все таблицы проектирую в 3НФ, а когда БД уже разрастется и,
--    если возникнут тормоза с JOIN, тогда уже денормализую
--    тут бы я сделал отдельную таблицу services, а тут писал бы UUID сервиса
    service_name TEXT    NOT NULL,
    price        INTEGER NOT NULL CHECK (price >= 0),
    user_id      UUID    NOT NULL,
    start_date   DATE    NOT NULL,
    end_date     DATE,
    created_at   TIMESTAMP        DEFAULT NOW(),
    updated_at   TIMESTAMP        DEFAULT NOW()
);

-- Вообще, при создании сервиса с нуля, я не ставлю индексы, т.к. на малозаполненной БД они будут только мешать
-- Потом, когда появятся медленные запросы, анализирую, и ствлю нужный индекс
-- Сейчас пишу, чтоб показать, что я знаю о них
CREATE INDEX idx_subscriptions_user_id
    ON subscriptions (user_id);

CREATE INDEX idx_subscriptions_service_name
    ON subscriptions (service_name);

CREATE INDEX idx_subscriptions_start_date ON subscriptions (start_date);
CREATE INDEX idx_subscriptions_end_date ON subscriptions (end_date);