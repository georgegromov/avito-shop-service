-- users
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- wallets (one per user)
CREATE TABLE wallets (
  user_id UUID PRIMARY KEY,
  balance BIGINT NOT NULL CHECK (balance >= 0),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),

	CONSTRAINT fk_wallet_owner_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- merch items (static catalog)
CREATE TABLE items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  price BIGINT NOT NULL CHECK (price >= 0)
);

-- purchases (when user buys an item)
CREATE TABLE purchases (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID,
  item_id UUID,
  quantity INT NOT NULL CHECK (quantity > 0),
  total_price BIGINT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),

	CONSTRAINT fk_buyer_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
	CONSTRAINT fk_item_id FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE SET NULL
);

-- coin transfers history (who sent to whom)
CREATE TABLE transfers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  from_user_id UUID,
  to_user_id UUID,
  amount BIGINT NOT NULL CHECK (amount > 0),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),

	CONSTRAINT fk_sender_user_id FOREIGN KEY (from_user_id) REFERENCES users(id) ON DELETE SET NULL,
	CONSTRAINT fk_receiver_user_id FOREIGN KEY (to_user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- indexes for read patterns
CREATE INDEX idx_purchases_user_id ON purchases(user_id);
CREATE INDEX idx_transfers_to_user_id ON transfers(to_user_id);
CREATE INDEX idx_transfers_from_user_id ON transfers(from_user_id);
