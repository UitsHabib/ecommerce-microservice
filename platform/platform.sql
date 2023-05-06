CREATE TYPE "user_status" AS ENUM (
  'active',
  'inactive'
);

CREATE TYPE "profile_type" AS ENUM (
  'custom',
  'standard'
);

CREATE TYPE "permission_type" AS ENUM (
  'custom',
  'standard'
);

CREATE TABLE "users" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "profile_id" UUID,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "first_name" varchar,
  "last_name" varchar,
  "phone" varchar,
  "status" user_status DEFAULT 'active',
  "last_login" timestamptz,
  "password_expiry_date" timestamptz,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_by" UUID,
  "updated_by" UUID,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "profiles" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "title" varchar UNIQUE NOT NULL,
  "slug" varchar UNIQUE NOT NULL,
  "type" profile_type DEFAULT 'custom',
  "description" varchar,
  "created_by" UUID,
  "updated_by" UUID,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "permissions" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "title" varchar UNIQUE NOT NULL,
  "slug" varchar UNIQUE NOT NULL,
  "type" permission_type DEFAULT 'custom',
  "description" varchar,
  "created_by" UUID,
  "updated_by" UUID,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "features" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "title" varchar UNIQUE NOT NULL,
  "slug" varchar UNIQUE NOT NULL,
  "parent_id" UUID,
  "description" varchar,
  "created_by" UUID,
  "updated_by" UUID,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "permission_features" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "permission_id" UUID,
  "feature_id" UUID,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "profile_permissions" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "permission_id" UUID,
  "profile_id" UUID,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "users" ADD FOREIGN KEY ("profile_id") REFERENCES "profiles" ("id");

ALTER TABLE "permission_features" ADD FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id");

ALTER TABLE "permission_features" ADD FOREIGN KEY ("feature_id") REFERENCES "features" ("id");

ALTER TABLE "profile_permissions" ADD FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id");

ALTER TABLE "profile_permissions" ADD FOREIGN KEY ("profile_id") REFERENCES "profiles" ("id");
