CREATE TYPE "comment_type" AS ENUM (
  'post',
  'merchant'
);

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
  "id" BIGSERIAL PRIMARY KEY,
  "owner" varchar NOT NULL,
  "image_url" varchar NOT NULL DEFAULT '/default/user/avatar.jpg',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "merchants" (
  "id" BIGSERIAL PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL DEFAULT 0,
  "profession" varchar NOT NULL,
  "title" varchar NOT NULL,
  "about" varchar NOT NULL,
  "image_url" varchar NOT NULL DEFAULT '/default/merchant/avatar.jpg',
  "rating" double PRECISION NOT NULL DEFAULT 0.0,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "posts" (
  "id" BIGSERIAL PRIMARY KEY,
  "merchant_id" bigint NOT NULL,
  "title" varchar,
  "image_url" varchar,
  "likes" int NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "comments" (
  "id" BIGSERIAL PRIMARY KEY,
  "comment_type" comment_type NOT NULL,
  "post_id" bigint,
  "merchant_id" bigint,
  "owner" varchar NOT NULL,
  "comment" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "consultancies" (
  "id" BIGSERIAL PRIMARY KEY,
  "merchant_id" bigint NOT NULL,
  "customer_id" bigint NOT NULL,
  "cost" bigint NOT NULL,
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

CREATE INDEX ON "merchants" ("owner");

CREATE INDEX ON "posts" ("merchant_id");

CREATE INDEX ON "consultancies" ("merchant_id");

CREATE INDEX ON "consultancies" ("customer_id");

CREATE INDEX ON "consultancies" ("merchant_id", "customer_id");

COMMENT ON COLUMN "posts"."title" IS 'can be null only if image_url is not null';

COMMENT ON COLUMN "posts"."image_url" IS 'can be null only if title is not null';

COMMENT ON COLUMN "comments"."post_id" IS 'cannot be null if comment_type is post';

COMMENT ON COLUMN "comments"."merchant_id" IS 'cannot be null if comment_type is merchant';

COMMENT ON COLUMN "consultancies"."cost" IS 'must be positive';

ALTER TABLE "comments" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("merchant_id") REFERENCES "merchants" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "customers" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "merchants" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "consultancies" ADD FOREIGN KEY ("merchant_id") REFERENCES "merchants" ("id");

ALTER TABLE "consultancies" ADD FOREIGN KEY ("customer_id") REFERENCES "customers" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "posts" ADD FOREIGN KEY ("merchant_id") REFERENCES "merchants" ("id");
