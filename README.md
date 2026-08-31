# Consultation Platform

A modular, production-oriented Go consultation platform built using **Domain-Driven Design (DDD)** and **Hexagonal / Clean Architecture**.

---

## 🏛 Architecture Overview

The system is organized into decoupled bounded contexts, communicating across defined boundaries via in-app interfaces, DTOs, and ports/adapters without tight circular dependencies:

```text
                       ┌─────────────────────────────────────┐
                       │             HTTP / API              │
                       └──────────────────┬──────────────────┘
                                          │
       ┌──────────────────┬───────────────┴───────────────┬──────────────────┐
       │                  │                               │                  │
       ▼                  ▼                               ▼                  ▼
┌──────────────┐   ┌──────────────┐               ┌───────────────┐   ┌──────────────┐
│   Identity   │   │  Consultant  │               │ Consultation  │   │    Expert    │
│    Module    │   │    Module    │               │  Case Module  │   │   Matching   │
└──────────────┘   └──────────────┘               └───────────────┘   └──────────────┘
```

---

## 🧠 Expert Matching Subsystem

Expert Matching is an in-app capability designed to match consultation requests with optimal consultant candidates through a deterministic filter followed by a swappable, provider-independent ranking stage (e.g. Rule-Based, ML, or AI/LLM).

### Matching Pipeline

```text
ConsultationCase
      │
      ▼
Candidate Generation  ──► Bounded by Category, Accepting Clients, Verification
      │
      ▼
CandidatePool         ──► Bounded slice (Max capacity, unique consultants)
      │
      ▼
Candidate Ranking     ──► Swappable Ranking Engine (Rules, ML, Gemini/AI)
      │
      ▼
Ranked Candidates     ──► Contiguous Ranks (1..N) & Monotonic Match Scores
      │
      ▼
Top N Projection      ──► Top 5 / Custom Slicing
```

---

### Core Domain Model (`internal/expertmatching/domain/`)

* **`MatchingRun` (Aggregate Root)**: Owns the full lifecycle (`PENDING` $\to$ `GENERATING` $\to$ `RANKING` $\to$ `COMPLETED` / `FAILED` / `CANCELLED`). Enforces candidate uniqueness, contiguous rank sequences ($1, 2, \dots, N$), and score monotonicity (no rank-score inversions).
* **`MatchingCategory` (Value Object)**: Strictly validated category representation that prevents arbitrary free-form strings from bypassing domain invariants.
* **`Expertise` (Value Object)**: Discrete skill / specialty tags (e.g., `"Distributed Systems"`, `"Go"`, `"PostgreSQL"`), decoupled from broad categories.
* **`CandidateProfile` (Value Object)**: Matching-specific profile containing only attributes relevant to ranking (`ConsultantID`, `Category`, `Profession`, `Expertise`, `YearsExperience`, `Bio`). Operational gating flags (`is_accepting_clients`, `is_verified`) are evaluated earlier by the candidate generator.
* **`CandidatePool` (Value Object)**: A bounded, deduplicated collection of candidates prepared for the ranking engine (enforces max pool limits).
* **`CandidateGenerationCriteria` (Value Object)**: Explicit search/filter parameters passed across module boundaries.
* **`RankedCandidate` (Value Object)**: An evaluated candidate containing an immutable rank, score ($0.0 \dots 1.0$), and structured explainability reasons.

---

## 👨‍💼 Consultant Subsystem

Manages consultant profiles, professions, availability slots, verifications, and discrete expertises.

### Consultant Expertise (`consultant_expertises`)
* Stored in a normalized relational table `consultant_expertises` (not arbitrary JSON).
* Supports complete CRUD operations:
  * `POST /api/v1/consultants/register`: Registers a consultant profile with optional initial skills.
  * `GET /api/v1/consultants/:id`: Retrieves public consultant profile with their expertises.
  * `GET /api/v1/consultants/me/expertises`: Lists the authenticated consultant's skills.
  * `POST /api/v1/consultants/me/expertises`: Adds a skill.
  * `PUT /api/v1/consultants/me/expertises`: Replaces skills in bulk.
  * `DELETE /api/v1/consultants/me/expertises/:expertiseID`: Removes a skill.

---

## ⚡ Asynchronous Architecture & Job Queue

Matching operations are triggered asynchronously:

1. **`POST /api/v1/consultation-cases/:id/match`**:
   - Creates a `PENDING` `MatchingRun` record in the database.
   - Enqueues a job payload (`{"run_id": "...", "case_id": "..."}`) into Redis (`queue:expertmatching:jobs`).
   - Responds immediately with `202 Accepted`.
2. **Background Worker / Execution**:
   - Picks up the job and invokes `ExecuteMatchingUsecase`.
   - Executes candidate generation, candidate ranking, and marks the run `COMPLETED` (or `FAILED` with details).
3. **`GET /api/v1/consultation-cases/:id/matches?top_n=5`**:
   - Queries the latest completed run and returns the top $N$ ranked candidates with structured explanation factors.

---

## 🧪 Testing

```bash
# Run all unit tests across all bounded contexts
go test -v ./internal/...

# Run expert matching and consultant usecase suites
go test -v ./internal/expertmatching/usecase/... ./internal/consultant/usecase/...
```
