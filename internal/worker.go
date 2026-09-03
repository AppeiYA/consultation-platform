package app

import (
	"context"
	"fmt"

	consultant_adapter "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/outbound/external/consultant"
	consultationcase_adapter "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/outbound/external/consultationcase"
	expertMatchingRepo "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/outbound/postgres/expertmatching"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/outbound/ranker"
	geminiranker "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/outbound/ranker/gemini"
	matching_worker "github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/worker"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/usecase"
	system "github.com/AppeiYA/consultation-platform/internal/shared/adapters/outbound/clock"
	"github.com/AppeiYA/consultation-platform/internal/shared/config"
	"github.com/AppeiYA/consultation-platform/internal/shared/db"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type Worker struct {
	config *config.Config
	server *asynq.Server
	mux    *asynq.ServeMux
	db     *db.DB
}

func NewWorker() (*Worker, error) {
	cfg := config.SetupConfig()
	logger.Init(cfg)

	database, err := db.Connect(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	repository := db.NewRepository(database)
	clock := system.NewSystemClock()

	// Setup ranker: plug-and-play Gemini ranker if API key is provided, otherwise fallback to RuleBasedRanker
	var candidateRanker outbound.CandidateRanker = ranker.NewRuleBasedRanker()
	if cfg.AI.GeminiAPIKey != "" {
		gRanker, err := geminiranker.New(context.Background(), cfg.AI.GeminiAPIKey, cfg.AI.GeminiModel)
		if err != nil {
			logger.Warn("failed to initialize gemini ranker, falling back to rule-based ranker", zap.Error(err))
		} else {
			candidateRanker = gRanker
			logger.Info("expert matching worker using Gemini ranker", zap.String("model", cfg.AI.GeminiModel))
		}
	} else {
		logger.Info("expert matching worker using rule-based ranker")
	}

	caseReader := consultationcase_adapter.NewCaseReaderAdapter(&repository)
	candidateGen := consultant_adapter.NewCandidateGeneratorAdapter(&repository)
	runRepo := expertMatchingRepo.NewExpertMatchingRepository(&repository)

	processMatching := usecase.NewProcessMatchingUsecase(
		caseReader,
		candidateGen,
		candidateRanker,
		runRepo,
		clock,
	)

	matchingWorker := matching_worker.NewMatchingWorker(processMatching)

	mux := asynq.NewServeMux()
	mux.HandleFunc(matching_worker.TypeStartMatching, matchingWorker.HandleStartMatching)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.Redis.Address,
			Username: cfg.Redis.Username,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"default": 1,
			},
		},
	)

	return &Worker{
		config: cfg,
		server: srv,
		mux:    mux,
		db:     database,
	}, nil
}

func (w *Worker) Run() error {
	defer func() {
		if w.db != nil {
			_ = w.db.Close()
		}
	}()

	logger.Info("expert matching background worker started listening for tasks")
	return w.server.Run(w.mux)
}
