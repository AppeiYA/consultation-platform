CREATE TABLE IF NOT EXISTS consultant_availabilities (
    id VARCHAR(64) PRIMARY KEY,
    consultant_id VARCHAR(64) NOT NULL,
    day_of_week SMALLINT NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT fk_availability_consultant
        FOREIGN KEY (consultant_id)
        REFERENCES consultants(id) ON DELETE CASCADE,

    CONSTRAINT chk_availability_day
        CHECK (day_of_week BETWEEN 0 AND 6),

    CONSTRAINT chk_availability_time
        CHECK (start_time < end_time),

    CONSTRAINT uq_consultant_day_time
        UNIQUE (consultant_id, day_of_week, start_time, end_time)
);

CREATE INDEX IF NOT EXISTS idx_consultant_availabilities_consultant_id
    ON consultant_availabilities (consultant_id);
CREATE INDEX IF NOT EXISTS idx_consultant_availabilities_day_of_week
    ON consultant_availabilities (day_of_week);
CREATE INDEX IF NOT EXISTS idx_consultant_availabilities_start_time
    ON consultant_availabilities (start_time);
CREATE INDEX IF NOT EXISTS idx_consultant_availabilities_end_time
    ON consultant_availabilities (end_time);