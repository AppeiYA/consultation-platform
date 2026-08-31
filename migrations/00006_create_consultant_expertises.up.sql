CREATE TABLE IF NOT EXISTS consultant_expertises (
    id VARCHAR(64) PRIMARY KEY,
    consultant_id VARCHAR(64) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_consultant_expertise
        FOREIGN KEY (consultant_id)
        REFERENCES consultants(id) ON DELETE CASCADE,

    CONSTRAINT uq_consultant_expertise
        UNIQUE (consultant_id, name)
);

CREATE INDEX IF NOT EXISTS idx_consultant_expertises_consultant_id ON consultant_expertises(consultant_id);
