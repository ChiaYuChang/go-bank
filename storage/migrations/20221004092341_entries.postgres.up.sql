CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" decimal NOT NULL DEFAULT 0,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "entries" ("account_id");

COMMENT ON COLUMN "entries"."amount" IS 'could be any number';

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");
