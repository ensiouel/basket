-- +goose Up
-- +goose StatementBegin
CREATE TABLE file
(
    id             UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    source_id      TEXT        NOT NULL,
    title          TEXT        NOT NULL,
    name           TEXT        NOT NULL,
    description    TEXT        NOT NULL DEFAULT '',
    size           BIGINT      NOT NULL,
    download_count INTEGER     NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL,
    updated_at     TIMESTAMPTZ NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE file;
-- +goose StatementEnd
