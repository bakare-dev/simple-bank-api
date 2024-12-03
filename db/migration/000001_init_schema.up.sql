-- Drop types if they already exist
DROP TYPE IF EXISTS "AccountType";
DROP TYPE IF EXISTS "TransactionType";
DROP TYPE IF EXISTS "TransactionStatus";

-- Create the types again
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

-- Create the tables
CREATE TABLE "Users" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" varchar(100) NOT NULL,
  "email" varchar(150) UNIQUE NOT NULL,
  "password" varchar(1000) NOT NULL,
  "phone_number" varchar(15) UNIQUE,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE "Accounts" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT NOT NULL, 
  "account_number" varchar(20) UNIQUE NOT NULL,
  "type" "AccountType" NOT NULL,
  "pin" varchar(1000) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE "Transactions" (
  "id" BIGSERIAL PRIMARY KEY,
  "account_id" BIGINT NOT NULL, 
  "type" "TransactionType" NOT NULL,
  "status" "TransactionStatus" NOT NULL,
  "amount" decimal(18,2) NOT NULL,
  "description" varchar(255),
  "transaction_date" timestamptz NOT NULL DEFAULT now()
);

-- Add foreign key constraints
ALTER TABLE "Accounts" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");
ALTER TABLE "Transactions" ADD FOREIGN KEY ("account_id") REFERENCES "Accounts" ("id");