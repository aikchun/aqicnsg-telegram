BEGIN;
ALTER TABLE subscription DROP COLUMN IF EXISTS "aqi_threshold";
COMMIT;
