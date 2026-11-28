# Skilloom Ingest

HTTP ingestion service for Skilloom. Accepts conversation events from `collector` hooks, authenticates tokens via Redis, and publishes normalized messages to NATS JetStream.

## Request Flow

![Ingest request flow](../assets/ingest-request-flow.png)

1. **Authenticate** — Extract the `Bearer` token, hash it with SHA-256, and resolve `(user_id, space_id)` from Redis (`token:<sha256>`).
2. **Validate** — Validate the incoming payload.
3. **Buffer** — Publish to NATS JetStream (`INGEST` stream, `ingest.messages`) with `event_id` as `Nats-Msg-Id` for deduplication.
4. **Acknowledge** — Return `202 Accepted` after JetStream confirms the publish.

## Configuration

Loaded from `config.yaml` (override path with `CONFIG_PATH`). Every setting can be overridden via environment variables:

| Key | Env Variable | Default |
|---|---|---|
| `http.addr` | `HTTP_ADDR` | `:8080` |
| `nats.url` | `NATS_URL` | `nats://127.0.0.1:4222` |
| `redis.addr` | `REDIS_ADDR` | `127.0.0.1:6379` |
| `redis.password` | `REDIS_PASSWORD` | `""` |
| `redis.db` | `REDIS_DB` | `0` |
