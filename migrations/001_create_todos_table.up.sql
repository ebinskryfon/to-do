CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE todos (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    completed   BOOLEAN NOT NULL DEFAULT FALSE,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    created_by  UUID,
    updated_by  UUID,
    deleted_by  UUID
);

CREATE INDEX idx_todos_is_active  ON todos(is_active);
CREATE INDEX idx_todos_created_at ON todos(created_at);
CREATE INDEX idx_todos_completed  ON todos(completed);
CREATE INDEX idx_todos_deleted_at ON todos(deleted_at);
