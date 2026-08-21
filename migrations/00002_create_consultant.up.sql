CREATE TABLE IF NOT EXISTS professions (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO professions (id, name, created_at) VALUES
    ('prof_9ee432d7-b672-40ae-b03f-c1f1fb696621', 'SOFTWARE_ENGINEER', CURRENT_TIMESTAMP),
    ('prof_12d965f5-e1f5-49aa-ac57-856772d236ce', 'LAWYER', CURRENT_TIMESTAMP),
    ('prof_940f840d-617a-4ead-8a8c-873851762bc7', 'DOCTOR', CURRENT_TIMESTAMP),
    ('prof_9a1e78f7-a053-4b4b-8802-fcac77e530ec', 'ACCOUNTANT', CURRENT_TIMESTAMP),
    ('prof_d95d5c58-d5be-4bca-bf87-b84b3f2a2681', 'THERAPIST', CURRENT_TIMESTAMP),
    ('prof_aef03e88-a0ee-49c9-b455-4b2210412b52', 'CLERGY', CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS consultants (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL UNIQUE,
    profession_id VARCHAR(64) NOT NULL,
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

    CONSTRAINT fk_consultant_profession
        FOREIGN KEY (profession_id)
        REFERENCES professions(id),

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