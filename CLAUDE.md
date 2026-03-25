# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

**Katerina** is a fully local voice AI English-club app. Users record English speech in the browser; the system transcribes it, corrects grammar, translates to Russian, generates an AI teacher response, and returns synthesised audio — all running inside Docker with no external API calls.

---

## Commands

### First-time setup
```bash
make tidy       # download Go dependencies (run once, outside Docker)
make setup      # build images, start stack, pull llama3 model (~4 GB)
```

### Daily workflow
```bash
make up         # start all containers (detached)
make down       # stop containers
make build      # rebuild images after code changes
make logs       # follow all logs
make restart    # restart all containers
make clean      # remove containers + volumes (destructive)
make pull-model # pull/update the Ollama LLM model
```

### Go (backend)
```bash
go build ./...                        # compile check
go test ./...                         # all tests
go test ./internal/core/...           # domain/use-case tests only
go test -run TestProcessVoice ./...   # single test by name
go vet ./...                          # static analysis
```

### Frontend
```bash
cd frontend
npm install
npm run dev     # dev server on :5173 (proxies /api → localhost:8080)
npm run build   # production build
```

### Services (local Python dev)
```bash
cd services/whisper && pip install -r requirements.txt && uvicorn app:app --port 8001
cd services/tts     && pip install -r requirements.txt && uvicorn app:app --port 8002
```

---

## Architecture

### Request flow
```
Browser (MediaRecorder)
  → POST /api/voice/process  [multipart audio]
  → Go backend
      1. ffmpeg  → WAV 16 kHz mono (temp file)
      2. Whisper → original text
      3. Ollama  → { corrected, translation, answer }  (JSON)
      4. Piper   → answer.wav  (saved to /app/audio volume)
      5. Postgres → save Conversation + Audio rows
  → JSON response { original, corrected, translation, answer, audio_url }
  → GET /api/audio/{id}  → serve WAV file
```

### Go backend — strict DDD / Clean Architecture
```
cmd/app/
  main.go              — server lifecycle (graceful shutdown)
  composition_root.go  — dependency injection wiring
  config.go            — YAML config structs + LoadConfig

internal/core/          ← NO infrastructure imports allowed here
  domain/model/         — User, Conversation, Audio (pure structs)
  ports/                — interfaces: repositories + STT/LLM/TTS ports
  application/usecases/ — ProcessVoiceUseCase, AuthUseCase, AudioUseCase

internal/adapters/
  in/http/              — chi router (package httpserver), handlers, middleware
  out/postgres/         — UserRepository, ConversationRepository, AudioRepository
  out/redis/            — TokenRepository (refresh tokens + blacklist)
  out/stt_client/       — WhisperClient  → POST /transcribe
  out/llm_client/       — OllamaClient   → POST /api/generate
  out/tts_client/       — PiperClient    → POST /synthesize
```

**Dependency rule:** `core/` never imports from `adapters/`. Adapters depend on ports (interfaces), not concrete types.

### AI services
| Service | Port | Image | Notes |
|---------|------|-------|-------|
| Whisper | 8001 | custom Python | `faster-whisper`, model downloaded at first transcription |
| Ollama  | 11434 | `ollama/ollama` | run `make pull-model` to pull `llama3` |
| Piper   | 8002 | custom Python | `en_US-lessac-medium` model baked into image |

### Auth
- Access token: JWT HS256, 15 min TTL
- Refresh token: UUID stored in Redis with 7-day TTL, one-time use (old token blacklisted on rotation)
- Demo credentials seeded at startup: `demo@katerina.local` / `password123`

### Storage
- Audio files: `/app/audio` volume (Docker), served via `GET /api/audio/{id}`
- Future: swap `PiperClient.audioStorage` path for MinIO presigned URLs

---

## Key config

`configs/config.yaml` — DB, Redis, JWT secret, AI service URLs, audio path.
Override at runtime via `CONFIG_PATH` env var.

---

## Extending

- **WebSocket streaming**: add a `/api/voice/stream` route; use `chi` SSE or upgrade to WebSocket in the handler layer only — use cases stay unchanged.
- **Kafka**: replace direct use-case calls with an event publisher behind a port interface.
- **MinIO**: implement a new `StoragePort` and inject it into `ProcessVoiceUseCase`.
- **Split microservices**: each use case maps cleanly to a separate service; ports already define the contracts.
