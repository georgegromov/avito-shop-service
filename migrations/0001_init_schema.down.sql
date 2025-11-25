DROP INDEX IF EXISTS idx_transfers_from_user_id;
DROP INDEX IF EXISTS idx_transfers_to_user_id;
DROP INDEX IF EXISTS idx_purchases_user_id;

DROP TABLE IF EXISTS transfers;
DROP TABLE IF EXISTS purchases;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS users;
