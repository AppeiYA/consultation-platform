CREATE TABLE IF NOT EXISTS consultants (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL UNIQUE,
    profession VARCHAR(100) NOT NULL,
    display_name VARCHAR(150) NOT NULL,
    bio TEXT NOT NULL,
    years_experience INTEGER NOT NULL,
    -- is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_accepting_clients BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT fk_consultant_user
        FOREIGN KEY (user_id)
        REFERENCES users(id),

    CONSTRAINT chk_consultant_years_experience
        CHECK (years_experience >= 0)
);

CREATE TABLE IF NOT EXISTS consultant_verifications (
    id VARCHAR(64) PRIMARY KEY,
    consultant_id VARCHAR(64) NOT NULL,
    provider VARCHAR(100) NOT NULL,
    provider_reference TEXT NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL CHECK (
    status IN (
        'NOT_STARTED',
        'PENDING',
        'REVIEW',
        'APPROVED',
        'REJECTED'
        )
    ),
    submitted_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT fk_verification_consultant
        FOREIGN KEY (consultant_id)
        REFERENCES consultants(id)
);