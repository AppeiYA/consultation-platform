CREATE TABLE IF NOT EXISTS consultation_cases (
    id VARCHAR(64) PRIMARY KEY,
    client_id VARCHAR(64) NOT NULL,
    title VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    category VARCHAR(100) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'DRAFT',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT fk_case_client
        FOREIGN KEY (client_id)
        REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_consultation_cases_client_id
    ON consultation_cases (client_id);
CREATE INDEX IF NOT EXISTS idx_consultation_cases_status
    ON consultation_cases (status);
CREATE INDEX IF NOT EXISTS idx_consultation_cases_category
    ON consultation_cases (category);
