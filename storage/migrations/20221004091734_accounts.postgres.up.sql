CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar(30) NOT NULL,
  "balance" decimal NOT NULL DEFAULT 0,
  "currency" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "deleted_at" timestamp
);

CREATE INDEX ON "accounts" ("owner");

