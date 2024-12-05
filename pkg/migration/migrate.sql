DROP TYPE IF EXISTS account_type;
DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS transaction_status;
CREATE TYPE account_type AS ENUM ('savings', 'current');
CREATE TYPE transaction_type AS ENUM ('credit', 'debit');
CREATE TYPE transaction_status AS ENUM ('processing', 'successful', 'failed');