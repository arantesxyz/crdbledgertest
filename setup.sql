CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    allow_negative BOOLEAN,
    balance NUMERIC(12, 2),
    CONSTRAINT check_account_allow_negative CHECK (
        allow_negative :: int between 0
        and 1
    ),
    CONSTRAINT check_account_positive_balance CHECK (balance * abs(allow_negative :: int - 1) >= 0)
);

CREATE TABLE transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID,
    amount NUMERIC(12, 2),
    isCredit BOOLEAN,
    status VARCHAR(10)
)