CREATE TABLE "apis" (
  "id" bigserial PRIMARY KEY,
  "key" varchar UNIQUE NOT NULL,
  "usage_count" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

COMMENT ON COLUMN "apis"."key" IS 'public keys from 3rd API';
