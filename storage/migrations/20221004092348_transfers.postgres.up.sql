CREATE TYPE "tstatus" AS ENUM (
  'success',
  'failure'
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "src_id" bigint NOT NULL,
  "dst_id" bigint NOT NULL,
  "amount" decimal NOT NULL DEFAULT 0,
  "status" tstatus NOT NULL DEFAULT 'failure',
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "transfers" ("src_id");

CREATE INDEX ON "transfers" ("dst_id");

CREATE INDEX ON "transfers" ("src_id", "dst_id");

COMMENT ON COLUMN "transfers"."amount" IS 'must gte 0';

ALTER TABLE "transfers" ADD FOREIGN KEY ("src_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("dst_id") REFERENCES "accounts" ("id");
