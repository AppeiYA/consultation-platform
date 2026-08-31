CREATE TABLE IF NOT EXISTS matching_runs (
    id VARCHAR(255) PRIMARY KEY,
    case_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    ranking_version VARCHAR(50) NOT NULL,
    failure_reason TEXT,
    cancellation_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (case_id) REFERENCES consultation_cases(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_matching_runs_case_id ON matching_runs(case_id);
CREATE INDEX IF NOT EXISTS idx_matching_runs_status ON matching_runs(status);

CREATE TABLE IF NOT EXISTS matching_run_candidates (
    id SERIAL PRIMARY KEY,
    run_id VARCHAR(255) NOT NULL,
    consultant_id VARCHAR(255) NOT NULL,
    rank_position INT NOT NULL,
    score NUMERIC(5,4) NOT NULL,
    reasons JSONB NOT NULL DEFAULT '[]',
    FOREIGN KEY (run_id) REFERENCES matching_runs(id) ON DELETE CASCADE,
    FOREIGN KEY (consultant_id) REFERENCES consultants(id) ON DELETE CASCADE,
    CONSTRAINT uq_run_consultant UNIQUE (run_id, consultant_id),
    CONSTRAINT uq_run_rank UNIQUE (run_id, rank_position)
);

CREATE INDEX IF NOT EXISTS idx_matching_run_candidates_run_id ON matching_run_candidates(run_id);
