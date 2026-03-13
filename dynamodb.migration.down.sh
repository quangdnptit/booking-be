#!/bin/bash

AWS_ACCESS_KEY_ID=fake
AWS_SECRET_ACCESS_KEY=fake
AWS_DEFAULT_REGION=us-east-1
ENDPOINT="--endpoint-url http://localhost:8000"
REGION="--region us-east-1"

set -e

echo "Deleting USERS table..."
aws dynamodb delete-table \
  --table-name users \
  $ENDPOINT $REGION || true

echo "Deleting BOOKINGS table..."
aws dynamodb delete-table \
  --table-name bookings \
  $ENDPOINT $REGION || true

echo "Deleting BOOKED_SEATS table..."
aws dynamodb delete-table \
  --table-name booked_seats \
  $ENDPOINT $REGION || true

echo "✅ All tables deleted"