CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TYPE IF EXISTS "AccountType";
DROP TYPE IF EXISTS "TransactionType";
DROP TYPE IF EXISTS "TransactionStatus";

CREATE TYPE "AccountType" AS ENUM (
  'Savings',
  'Checking'
);

CREATE TYPE "TransactionType" AS ENUM (
  'Debit',
  'Credit'
);

CREATE TYPE "TransactionStatus" AS ENUM (
  'Processing',
  'Successful',
  'Failed'
);

CREATE TABLE "Users" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "name" varchar(100) NOT NULL,
  "email" varchar(150) UNIQUE NOT NULL,
  "password" varchar(1000) NOT NULL,
  "phone_number" varchar(15) UNIQUE,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE "Accounts" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "user_id" UUID NOT NULL, 
  "account_number" varchar(20) UNIQUE NOT NULL,
  "type" "AccountType" NOT NULL,
  "pin" varchar(1000) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE "Transactions" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "account_id" UUID NOT NULL, 
  "type" "TransactionType" NOT NULL,
  "status" "TransactionStatus" NOT NULL,
  "amount" decimal(18,2) NOT NULL,
  "description" varchar(255),
  "transaction_date" timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE "Accounts" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");
ALTER TABLE "Transactions" ADD FOREIGN KEY ("account_id") REFERENCES "Accounts" ("id");