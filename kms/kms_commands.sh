#!/bin/bash

TRANSACTION_URL="${TRANSACTION_URL:-https://apis.greyboxpay.com/api/v1/transaction/verify}"
PERIOD_OF_CHECK="${PERIOD_OF_CHECK:-300}"
WALLET_STORAGE_LOCATION="${WALLET_STORAGE_LOCATION:-$(pwd)/wallet/wallet.dat}"
CHAINS="${CHAINS:-CELO,XLM}"
ENVFILE="${ENVFILE:-$(pwd)/.env}"

# Load environment variables from the specified .env file

# Start the Tatum KMS daemon with specified configurations
exec docker run -it --env-file .env -v $pwd:/root/.tatumrc tatumio/tatum-kms daemon \
  --external-url="$TRANSACTION_URL" \
  --period="$PERIOD_OF_CHECK" \
  --chain="$CHAINS" \
  --env-file="$ENVFILE"

