CREATE TABLE IF NOT EXISTS launchpads (
    id bigserial PRIMARY KEY,
    chain_id bigint NOT NULL,
    launchpad_id bytea NOT NULL,
    slug text NOT NULL,
    name text NOT NULL,
    description text NOT NULL DEFAULT '',
    logo_url text NOT NULL DEFAULT '',
    primary_color text NOT NULL DEFAULT '',
    owner bytea NOT NULL,
    treasury bytea NOT NULL,
    registry_tx_hash bytea NOT NULL,
    registry_block_number bigint NOT NULL,
    active boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT launchpads_chain_slug_key UNIQUE (chain_id, slug),
    CONSTRAINT launchpads_chain_launchpad_id_key UNIQUE (chain_id, launchpad_id)
);

CREATE TABLE IF NOT EXISTS token_launches (
    id bigserial PRIMARY KEY,
    launchpad_pk bigint NOT NULL REFERENCES launchpads(id),
    chain_id bigint NOT NULL,
    token bytea NOT NULL,
    pool_id bytea NOT NULL,
    creator bytea NOT NULL,
    quote bytea NOT NULL,
    factory bytea NOT NULL,
    supply numeric(78, 0) NOT NULL,
    tx_hash bytea NOT NULL,
    block_number bigint NOT NULL,
    block_hash bytea NOT NULL,
    log_index integer NOT NULL,
    status text NOT NULL,
    created_at timestamptz NOT NULL,
    CONSTRAINT token_launches_chain_token_key UNIQUE (chain_id, token),
    CONSTRAINT token_launches_chain_pool_key UNIQUE (chain_id, pool_id),
    CONSTRAINT token_launches_chain_tx_log_key UNIQUE (chain_id, tx_hash, log_index)
);

CREATE INDEX IF NOT EXISTS token_launches_launchpad_idx
    ON token_launches (launchpad_pk, block_number DESC, log_index DESC);

CREATE TABLE IF NOT EXISTS chain_transactions (
    id bigserial PRIMARY KEY,
    chain_id bigint NOT NULL,
    tx_hash bytea NOT NULL,
    kind text NOT NULL,
    status text NOT NULL,
    block_number bigint,
    block_hash bytea,
    failure_reason text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT chain_transactions_chain_tx_hash_key UNIQUE (chain_id, tx_hash)
);

CREATE TABLE IF NOT EXISTS indexer_checkpoints (
    chain_id bigint NOT NULL,
    factory bytea NOT NULL,
    next_block bigint NOT NULL,
    last_block_hash bytea,
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (chain_id, factory)
);

CREATE TABLE IF NOT EXISTS auth_nonces (
    id uuid PRIMARY KEY,
    chain_id bigint NOT NULL,
    wallet bytea NOT NULL,
    nonce_hash bytea NOT NULL UNIQUE,
    message text NOT NULL,
    expires_at timestamptz NOT NULL,
    used_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS auth_nonces_wallet_expiry_idx
    ON auth_nonces (chain_id, wallet, expires_at DESC);
