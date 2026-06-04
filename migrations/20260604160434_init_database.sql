-- +goose Up
-- +goose StatementBegin
CREATE DATABASE analytics ON CLUSTER study_cluster;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP DATABASE analytics;
-- +goose StatementEnd
