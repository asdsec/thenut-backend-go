CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "phone_number" varchar UNIQUE NOT NULL,
  "image_url" varchar NOT NULL DEFAULT '/default/user/avatar.jpg',
  "gender" varchar(1) NOT NULL DEFAULT 'm',
  "disabled" boolean NOT NULL DEFAULT false,
  "birth_date" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "customers" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "image_url" varchar NOT NULL DEFAULT '/default/user/avatar.jpg',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "merchants" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "profession" varchar NOT NULL,
  "title" varchar NOT NULL,
  "about" varchar NOT NULL,
  "image_url" varchar NOT NULL DEFAULT '/default/merchant/avatar.jpg',
  "rating" integer NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "consultancies" (
  "id" bigserial PRIMARY KEY,
  "merchant_id" bigint NOT NULL,
  "customer_id" bigint NOT NULL,
  "cost" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "customer_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "type" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "payments" (
  "id" bigserial PRIMARY KEY,
  "merchant_id" bigint NOT NULL,
  "customer_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- CREATE UNIQUE INDEX ON "customers" ("owner");
ALTER TABLE "customers" ADD CONSTRAINT "owner_key" UNIQUE ("owner");

CREATE INDEX ON "merchants" ("owner");

-- CREATE UNIQUE INDEX ON "merchants" ("owner", "profession");
ALTER TABLE "merchants" ADD CONSTRAINT "owner_profession_key" UNIQUE ("owner", "profession");

CREATE INDEX ON "consultancies" ("merchant_id");

CREATE INDEX ON "consultancies" ("customer_id");

CREATE INDEX ON "consultancies" ("merchant_id", "customer_id");

CREATE INDEX ON "entries" ("customer_id");

CREATE INDEX ON "payments" ("merchant_id");

CREATE INDEX ON "payments" ("customer_id");

CREATE INDEX ON "payments" ("merchant_id", "customer_id");

COMMENT ON COLUMN "consultancies"."cost" IS 'must be positive';

COMMENT ON COLUMN "entries"."amount" IS 'might be positive or negative';

COMMENT ON COLUMN "payments"."amount" IS 'must be positive';

ALTER TABLE "customers" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "merchants" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "consultancies" ADD FOREIGN KEY ("merchant_id") REFERENCES "merchants" ("id");

ALTER TABLE "consultancies" ADD FOREIGN KEY ("customer_id") REFERENCES "customers" ("id");

ALTER TABLE "entries" ADD FOREIGN KEY ("customer_id") REFERENCES "customers" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("merchant_id") REFERENCES "merchants" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("customer_id") REFERENCES "customers" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

-- ALTER TABLE "entries" ADD FOREIGN KEY ("created_at") REFERENCES "entries" ("id");
