CREATE TABLE IF NOT EXISTS consultants (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL UNIQUE,
    profession VARCHAR(100) NOT NULL,
    display_name VARCHAR(150) NOT NULL,
    bio TEXT NOT NULL,
    years_experience INTEGER NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_accepting_clients BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT fk_consultant_user
        FOREIGN KEY (user_id)
        REFERENCES users(id),

    CONSTRAINT chk_consultant_years_experience
        CHECK (years_experience >= 0)
);