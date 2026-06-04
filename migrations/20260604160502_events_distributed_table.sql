-- Это НЕ хранение данных.
--
-- Это:
--
-- маршрутизатор в кластер

-- +goose Up
-- +goose StatementBegin
CREATE TABLE analytics.events_all ON CLUSTER study_cluster
(
    id UInt64,
    user_id UInt64,
    event_name String,
    created_at DateTime
)
    ENGINE = Distributed(
    'study_cluster', -- из remote_server в config/yaml
    'analytics',
    'events',
    rand()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE analytics.events_all;
-- +goose StatementEnd
