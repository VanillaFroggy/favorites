CREATE TABLE IF NOT EXISTS favorites
(
    id          UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    project_id  UUID             NOT NULL,
    owner_type  VARCHAR          NOT NULL,
    owner_id    UUID             NOT NULL,
    object_id   UUID             NOT NULL,
    object_type VARCHAR          NOT NULL,
    created_at  TIMESTAMP        NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_favorites_project_id ON favorites (project_id);
CREATE INDEX IF NOT EXISTS idx_favorites_owner ON favorites (owner_id, owner_type);
CREATE INDEX IF NOT EXISTS idx_favorites_object ON favorites (object_id, object_type);
CREATE INDEX IF NOT EXISTS idx_favorites_created_at ON favorites (created_at);