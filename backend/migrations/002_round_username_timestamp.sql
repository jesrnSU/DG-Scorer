BEGIN;

-- Earlier rounds have no saved username or time. Keep them as unknown at midnight UTC.
ALTER TABLE rounds ADD COLUMN username TEXT NOT NULL DEFAULT 'unknown';
ALTER TABLE rounds ALTER COLUMN username DROP DEFAULT;
ALTER TABLE rounds ALTER COLUMN date TYPE TIMESTAMPTZ
    USING date::timestamp AT TIME ZONE 'UTC';

COMMIT;
