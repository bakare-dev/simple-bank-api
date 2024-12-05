-- -------------------------------------------------------------
-- TablePlus 6.2.0(576)
--
-- https://tableplus.com/
--
-- Database: simplebank
-- Generation Time: 2024-12-05 05:31:42.0290
-- -------------------------------------------------------------






















DROP TABLE IF EXISTS "public"."accounts";
-- This script only contains the table creation statements and does not fully represent the table in the database. Do not use it as a backup.

DROP TYPE IF EXISTS "public"."account_type";
CREATE TYPE "public"."account_type" AS ENUM ('savings', 'current');

-- Table Definition
CREATE TABLE "public"."accounts" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "user_id" uuid NOT NULL,
    "number" varchar(15) NOT NULL,
    "type" "public"."account_type" NOT NULL,
    "pin" text NOT NULL,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."profiles";
-- This script only contains the table creation statements and does not fully represent the table in the database. Do not use it as a backup.

-- Table Definition
CREATE TABLE "public"."profiles" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "user_id" uuid NOT NULL,
    "first_name" varchar(255) NOT NULL,
    "last_name" varchar(255) NOT NULL,
    "phone_number" varchar(20),
    "date_of_birth" text,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."transactions";
-- This script only contains the table creation statements and does not fully represent the table in the database. Do not use it as a backup.

DROP TYPE IF EXISTS "public"."transaction_type";
CREATE TYPE "public"."transaction_type" AS ENUM ('credit', 'debit');
DROP TYPE IF EXISTS "public"."transaction_status";
CREATE TYPE "public"."transaction_status" AS ENUM ('processing', 'successful', 'failed');

-- Table Definition
CREATE TABLE "public"."transactions" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "account_id" uuid NOT NULL,
    "amount" numeric(20,2) NOT NULL,
    "type" "public"."transaction_type" NOT NULL,
    "status" "public"."transaction_status" NOT NULL,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."users";
-- This script only contains the table creation statements and does not fully represent the table in the database. Do not use it as a backup.

-- Table Definition
CREATE TABLE "public"."users" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "email" text NOT NULL,
    "password" text NOT NULL,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "role" text NOT NULL DEFAULT 'customer'::text,
    "status" text NOT NULL DEFAULT 'not_activated'::text,
    PRIMARY KEY ("id")
);

;
;
;
;
;
;
;
;
;
;
ALTER TABLE "public"."accounts" ADD FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE ON UPDATE CASCADE;


-- Indices
CREATE UNIQUE INDEX idx_accounts_user_id ON public.accounts USING btree (user_id);
CREATE UNIQUE INDEX idx_accounts_number ON public.accounts USING btree (number);
ALTER TABLE "public"."profiles" ADD FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE ON UPDATE CASCADE;


-- Indices
CREATE UNIQUE INDEX idx_profiles_user_id ON public.profiles USING btree (user_id);
ALTER TABLE "public"."transactions" ADD FOREIGN KEY ("account_id") REFERENCES "public"."accounts"("id") ON DELETE CASCADE ON UPDATE CASCADE;


-- Indices
CREATE UNIQUE INDEX uni_users_email ON public.users USING btree (email);
CREATE INDEX idx_user_email ON public.users USING btree (email);
