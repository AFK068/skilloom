# Skilloom Ingest

HTTP ingestion service for Skilloom. It accepts conversation events from `collector` hooks.

## Target Request Flow

The current HTTP handler validates the payload and returns `202 Accepted` after calling `SaveMessages`, which is still a stub. Redis token verification and JetStream publishing are planned.

![](../assets/ingest-request-flow.png)

1. **Authenticate** — Extract the `Bearer` token, hash it with SHA-256, and resolve `(user_id, space_id)` from Redis.
2. **Validate & Normalize** — Check the incoming payload and wrap it into a slice (`[]Message`) so single-message and future batch endpoints can share the same logic.
3. **Buffer Durably** — Publish to the file-backed NATS JetStream `INGEST` stream (`ingest.>`, `WorkQueuePolicy`) with `event_id` as `Nats-Msg-Id` for deduplication.
4. **Acknowledge** — Return `202 Accepted` after JetStream confirms the publish. The collector uses a 3-second HTTP timeout.

## Layout

```text
ingest/
├── cmd/ingest/                  # Application entrypoint & Uber Fx wiring
├── config.yaml                  # Default local configuration
└── internal/
    ├── config/                  # YAML + env loader (cleanenv)
    ├── domain/                  # Core models (Message, Principal)
    ├── natsclient/              # NATS JetStream client, stream setup & publish
    ├── services/messages/       # Ingestion service & repository/publisher ports
    ├── store/redis/             # Redis client & KV operations
    └── transport/http/          # Gin server, Fx route groups & /v1 handlers
```

## Configuration

Loaded from `config.yaml` (override path with `CONFIG_PATH`). Every setting can be overridden via environment variables:

| Key | Env Variable | Default |
|---|---|---|
| `http.addr` | `HTTP_ADDR` | `:8080` |
| `nats.url` | `NATS_URL` | `nats://127.0.0.1:4222` |
| `redis.addr` | `REDIS_ADDR` | `127.0.0.1:6379` |
| `redis.password` | `REDIS_PASSWORD` | `""` |
| `redis.db` | `REDIS_DB` | `0` |

