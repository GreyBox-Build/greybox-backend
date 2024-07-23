#!/bin/bash


TRANSACTION_URL="https://apis.greyboxpay.com/api/v1/transaction/verify/"
PERIOD_OF_CHECK=300
WALLET_STORAGE_LOCATION="/kms/wallet/wallet.dat"
CHAINS="CELO,XLM"

# Start the Tatum KMS daemon with specified configurations
docker run -d --name tatum-kms \
  tatumio/tatum-kms \
  daemon \
  --external-url=$TRANSACTION_URL \
  --period=$PERIOD_OF_CHECK \
  --path=$WALLET_STORAGE_LOCATION \
  --chain=$CHAINS
