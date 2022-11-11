CREATE TABLE "currencies" (
  "id" serial PRIMARY KEY,
  "name" varchar(50) NOT NULL,
  "abbr" varchar(5) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "deleted_at" timestamp
);

CREATE INDEX ON "currencies" ("abbr");

ALTER TABLE "accounts" ADD FOREIGN KEY ("currency") REFERENCES "currencies" ("id");
