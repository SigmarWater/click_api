-- +goose Up
-- +goose StatementBegin

-- 1. MergeTree — базовый движок ClickHouse
--
-- Он даёт:
--
-- быстрые вставки
-- быстрые SELECT
-- сортировку по ключу

-- ORDER BY — это СУПЕР важно
-- данные физически сортируются по времени
-- запросы по времени работают очень быстро
CREATE TABLE analytics.events ON CLUSTER study_cluster(
    id UInt64,
    user_id UInt64,
    event_name String,
    create_at DateTime
)
ENGINE = MergeTree()
ORDER BY (create_at, id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE analytics.events;
-- +goose StatementEnd
