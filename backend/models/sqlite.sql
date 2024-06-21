
UPDATE wallet_addresses
SET "wallet_index" = 50
WHERE public_address = '0x469c6aFBA896eA01400d2c79C2f34189E4DCCd00';

--INSERT INTO wallet_addresses (
    public_address,
    is_active,
    wallet_chain,
    master_wallet_id,
    in_use,
    created_at,
    updated_at,
    deleted_at
--) VALUES (
    '0x469c6aFBA896eA01400d2c79C2f34189E4DCCd00',
    true,
    'CELO',
    1,
    false,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    NULL
--);