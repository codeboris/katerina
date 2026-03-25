package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	httpserver "github.com/codeboris/katerina/internal/adapters/in/http"
	"github.com/codeboris/katerina/internal/adapters/out/llm_client"
	pgadapter "github.com/codeboris/katerina/internal/adapters/out/postgres"
	redisadapter "github.com/codeboris/katerina/internal/adapters/out/redis"
	"github.com/codeboris/katerina/internal/adapters/out/stt_client"
	"github.com/codeboris/katerina/internal/adapters/out/tts_client"
	"github.com/codeboris/katerina/internal/core/application/usecases"
	"github.com/go-chi/chi/v5"
)

type App struct {
	router *chi.Mux
	cfg    *Config
}

func NewApp(cfg *Config) (*App, error) {
	// ── Postgres ──────────────────────────────────────────────────────────────
	db, err := sql.Open("postgres", cfg.DB.DSN())
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	slog.Info("postgres connected")

	// ── Redis ─────────────────────────────────────────────────────────────────
	rdb := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr()})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	slog.Info("redis connected")

	// ── Repositories ──────────────────────────────────────────────────────────
	userRepo := pgadapter.NewUserRepository(db)
	convRepo := pgadapter.NewConversationRepository(db)
	audioRepo := pgadapter.NewAudioRepository(db)
	tokenRepo := redisadapter.NewTokenRepository(rdb)

	// ── AI Clients ────────────────────────────────────────────────────────────
	sttClient := stt_client.NewWhisperClient(cfg.AI.STTURL)
	llmClient := llm_client.NewOllamaClient(cfg.AI.LLMURL, cfg.AI.LLMModel)
	ttsClient := tts_client.NewPiperClient(cfg.AI.TTSURL, cfg.Storage.AudioPath)

	// ── Use Cases ─────────────────────────────────────────────────────────────
	authUC := usecases.NewAuthUseCase(userRepo, tokenRepo, usecases.AuthConfig{
		JWTSecret:  cfg.JWT.Secret,
		AccessTTL:  cfg.JWT.AccessTTL,
		RefreshTTL: cfg.JWT.RefreshTTL,
	})
	processUC := usecases.NewProcessVoiceUseCase(sttClient, llmClient, ttsClient, convRepo, audioRepo, cfg.Storage.AudioPath)
	audioUC := usecases.NewAudioUseCase(audioRepo)

	// Seed demo user
	if err := authUC.SeedMockUser(context.Background()); err != nil {
		slog.Warn("seed mock user", "error", err)
	}

	// ── HTTP Router ───────────────────────────────────────────────────────────
	router := httpserver.NewRouter(authUC, processUC, audioUC)

	return &App{router: router, cfg: cfg}, nil
}
