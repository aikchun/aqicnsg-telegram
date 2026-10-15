BEGIN;
CREATE TABLE IF NOT EXISTS subscription (
	"id" bigserial PRIMARY KEY,
	"chat_id" bigint NOT NULL UNIQUE,
	"is_active" BOOLEAN NOT NULL DEFAULT '1',
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
);
COMMIT;