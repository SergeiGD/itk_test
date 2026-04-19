-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    balance decimal NOT NULL,
    created_at timestamp NOT NULL default current_timestamp
    CONSTRAINT balance_positive CHECK (balance >= 0)
);
CREATE TABLE operations (
    id SERIAL PRIMARY KEY,
    wallet_id UUID references wallets (id) NOT NULL,
    type VARCHAR(255) NOT NULL,
    amount decimal NOT NULL,
    created_at timestamp NOT NULL default current_timestamp
    CONSTRAINT amount_positive CHECK (amount >= 0)
);
INSERT INTO wallets (id, balance, created_at) VALUES ('1d8a7307-b5f8-4686-b9dc-b752430abbd8', 0, now());
INSERT INTO wallets (id, balance, created_at) VALUES ('69359037-9599-48e7-b8f2-48393c019135', 0, now());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE wallets;
DROP TABLE operations;
-- +goose StatementEnd
