#!/bin/sh

TRANSACTION_URL="${TRANSACTION_URL:-https://apis.greyboxpay.com/api/v1/transaction/verify}"
PERIOD_OF_CHECK="${PERIOD_OF_CHECK:-300}"
WALLET_STORAGE_LOCATION="${WALLET_STORAGE_LOCATION:-/kms/wallet/wallet.dat}"
CHAINS="${CHAINS:-CELO,XLM}"

# Start the Tatum KMS daemon with specified configurations
exec tatum-kms daemon \
  --external-url="$TRANSACTION_URL" \
  --period="$PERIOD_OF_CHECK" \
  --path="$WALLET_STORAGE_LOCATION" \
  --chain="$CHAINS"

