CREATE TABLE consultant_availabilities (
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
        REFERENCES consultants(id),

    CONSTRAINT chk_availability_day
        CHECK (day_of_week BETWEEN 0 AND 6),

    CONSTRAINT chk_availability_time
        CHECK (start_time < end_time),

    CONSTRAINT uq_consultant_day_time
        UNIQUE (consultant_id, day_of_week, start_time, end_time)
);