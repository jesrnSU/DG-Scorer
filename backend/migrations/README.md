# Database migrations

Run migrations explicitly in order. There is no automatic migration runner.

For a new database, from the project root:

```sh
docker compose exec -T db psql -U postgres -d postgres -v ON_ERROR_STOP=1 < backend/migrations/001_create_rounds.sql
docker compose exec -T db psql -U postgres -d postgres -v ON_ERROR_STOP=1 < backend/migrations/002_round_username_timestamp.sql
```

If tables already exist, run only migration 002 once. It adds username and changes
date to TIMESTAMPTZ without removing rounds or hole scores. Existing rounds get
username `unknown` and midnight UTC because their original names and times were
not stored. New API requests require a username and RFC 3339 timestamp with timezone.
Timestamps are returned in UTC, preserving the instant from the CSV offset.
