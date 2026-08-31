package usecase

import (
	"math"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
	"github.com/stretchr/testify/require"
)

func TestMatchingCategory_Invariants(t *testing.T) {
	t.Run("should create valid category", func(t *testing.T) {
		c, err := domain.NewMatchingCategory("SOFTWARE_ENGINEERING")
		require.NoError(t, err)
		require.Equal(t, "SOFTWARE_ENGINEERING", c.Value())
		require.Equal(t, "SOFTWARE_ENGINEERING", c.String())

		c2, _ := domain.NewMatchingCategory("software_engineering")
		require.True(t, c.Equals(c2))
	})

	t.Run("should reject empty or whitespace category", func(t *testing.T) {
		_, errEmpty := domain.NewMatchingCategory("")
		require.Equal(t, domain.ErrInvalidMatchingCategory, errEmpty)

		_, errWhitespace := domain.NewMatchingCategory("   ")
		require.Equal(t, domain.ErrInvalidMatchingCategory, errWhitespace)
	})
}

func TestExpertise_Invariants(t *testing.T) {
	t.Run("should create valid expertise item", func(t *testing.T) {
		e, err := domain.NewExpertise("Distributed Systems")
		require.NoError(t, err)
		require.Equal(t, "Distributed Systems", e.Name())
		require.Equal(t, "Distributed Systems", e.String())

		e2, _ := domain.NewExpertise("distributed systems")
		require.True(t, e.Equals(e2))
	})

	t.Run("should reject empty expertise", func(t *testing.T) {
		_, err := domain.NewExpertise("   ")
		require.Equal(t, domain.ErrInvalidExpertise, err)
	})
}

func TestCandidateProfile_Invariants(t *testing.T) {
	cat, _ := domain.NewMatchingCategory("SOFTWARE_ENGINEERING")
	exp1, _ := domain.NewExpertise("Go")
	exp2, _ := domain.NewExpertise("PostgreSQL")

	t.Run("should create valid candidate profile with defensive copying", func(t *testing.T) {
		expertiseList := []domain.Expertise{exp1, exp2}
		cp, err := domain.NewCandidateProfile(
			"con_100",
			cat,
			"Software Engineer",
			expertiseList,
			7,
			"Senior engineer bio",
		)
		require.NoError(t, err)
		require.Equal(t, "con_100", cp.ConsultantID())
		require.Equal(t, cat, cp.Category())
		require.Equal(t, "Software Engineer", cp.Profession())
		require.Len(t, cp.Expertise(), 2)
		require.Equal(t, 7, cp.YearsExperience())
		require.Equal(t, "Senior engineer bio", cp.Bio())

		// Test defensive copying on egress
		ret := cp.Expertise()
		newExp, _ := domain.NewExpertise("React")
		ret[0] = newExp
		require.Equal(t, "Go", cp.Expertise()[0].Name())
	})

	t.Run("should reject negative years of experience", func(t *testing.T) {
		_, err := domain.NewCandidateProfile(
			"con_100",
			cat,
			"Software Engineer",
			nil,
			-1,
			"bio",
		)
		require.Equal(t, domain.ErrInvalidYearsExperience, err)
	})

	t.Run("should reject empty consultant ID", func(t *testing.T) {
		_, err := domain.NewCandidateProfile(
			"",
			cat,
			"Software Engineer",
			nil,
			5,
			"bio",
		)
		require.Equal(t, domain.ErrInvalidConsultantID, err)
	})
}

func TestCandidatePool_Invariants(t *testing.T) {
	cat, _ := domain.NewMatchingCategory("LEGAL")
	p1, _ := domain.NewCandidateProfile("con_1", cat, "Lawyer", nil, 5, "bio")
	p2, _ := domain.NewCandidateProfile("con_2", cat, "Lawyer", nil, 7, "bio")

	t.Run("should create valid candidate pool", func(t *testing.T) {
		pool, err := domain.NewCandidatePool([]domain.CandidateProfile{p1, p2}, 10)
		require.NoError(t, err)
		require.Equal(t, 2, pool.Size())
		require.False(t, pool.IsEmpty())
	})

	t.Run("should reject duplicate consultant in candidate pool", func(t *testing.T) {
		dup, _ := domain.NewCandidateProfile("con_1", cat, "Lawyer", nil, 3, "bio")
		_, err := domain.NewCandidatePool([]domain.CandidateProfile{p1, dup}, 10)
		require.Equal(t, domain.ErrDuplicateCandidateInPool, err)
	})

	t.Run("should reject candidate pool exceeding max capacity limit", func(t *testing.T) {
		_, err := domain.NewCandidatePool([]domain.CandidateProfile{p1, p2}, 1) // max capacity = 1, but 2 provided
		require.Equal(t, domain.ErrCandidatePoolExceeded, err)
	})
}

func TestMatchScore_Invariants(t *testing.T) {
	t.Run("valid range", func(t *testing.T) {
		s, err := domain.NewMatchScore(0.85)
		require.NoError(t, err)
		require.Equal(t, 0.85, s.Value())
	})

	t.Run("boundary checks", func(t *testing.T) {
		s0, err0 := domain.NewMatchScore(0.0)
		require.NoError(t, err0)
		require.Equal(t, 0.0, s0.Value())

		s1, err1 := domain.NewMatchScore(1.0)
		require.NoError(t, err1)
		require.Equal(t, 1.0, s1.Value())

		_, errNeg := domain.NewMatchScore(-0.01)
		require.Equal(t, domain.ErrInvalidMatchScore, errNeg)

		_, errOver := domain.NewMatchScore(1.001)
		require.Equal(t, domain.ErrInvalidMatchScore, errOver)

		_, errNaN := domain.NewMatchScore(math.NaN())
		require.Equal(t, domain.ErrInvalidMatchScore, errNaN)
	})
}

func TestMatchingRun_ContiguousRankAndScoreMonotonicity_Invariants(t *testing.T) {
	now := time.Date(2026, 8, 27, 12, 0, 0, 0, time.UTC)
	v1, _ := domain.NewRankingVersion("v1")

	t.Run("should reject non-contiguous rank sequence", func(t *testing.T) {
		run, _ := domain.NewMatchingRun("mrun_rank_gap", "case_1", v1, now)
		_ = run.StartGeneration()
		_ = run.StartRanking()

		r1, _ := domain.NewRank(1)
		r3, _ := domain.NewRank(3) // skipped rank 2!
		s1, _ := domain.NewMatchScore(0.9)
		s2, _ := domain.NewMatchScore(0.8)

		c1, _ := domain.NewRankedCandidate("con_1", r1, s1, nil)
		c2, _ := domain.NewRankedCandidate("con_2", r3, s2, nil)

		err := run.Complete([]domain.RankedCandidate{c1, c2}, now)
		require.Equal(t, domain.ErrNonContiguousRank, err)
	})

	t.Run("should reject contradictory rank vs score (lower rank with higher score)", func(t *testing.T) {
		run, _ := domain.NewMatchingRun("mrun_score_inversion", "case_1", v1, now)
		_ = run.StartGeneration()
		_ = run.StartRanking()

		r1, _ := domain.NewRank(1)
		r2, _ := domain.NewRank(2)
		sLow, _ := domain.NewMatchScore(0.5)
		sHigh, _ := domain.NewMatchScore(0.9) // rank 2 has 0.9 while rank 1 has 0.5 (contradiction!)

		c1, _ := domain.NewRankedCandidate("con_1", r1, sLow, nil)
		c2, _ := domain.NewRankedCandidate("con_2", r2, sHigh, nil)

		err := run.Complete([]domain.RankedCandidate{c1, c2}, now)
		require.Equal(t, domain.ErrContradictoryRankScore, err)
	})
}
