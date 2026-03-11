#!/bin/bash

AWS_ACCESS_KEY_ID=fake
AWS_SECRET_ACCESS_KEY=fake
AWS_DEFAULT_REGION=us-east-1
ENDPOINT="--endpoint-url http://localhost:8000"
REGION="--region us-east-1"

set -e  # stop script if any command fails

echo "Creating BOOKINGS table..."

aws dynamodb create-table \
  $ENDPOINT $REGION \
  --table-name bookings \
  --attribute-definitions \
      AttributeName=id,AttributeType=S \
      AttributeName=user_id,AttributeType=S \
      AttributeName=showtime_id,AttributeType=S \
      AttributeName=status,AttributeType=S \
      AttributeName=created_at,AttributeType=S \
  --key-schema \
      AttributeName=id,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --global-secondary-indexes '[
      {
        "IndexName": "user-bookings-index",
        "KeySchema": [
          {"AttributeName":"user_id","KeyType":"HASH"},
          {"AttributeName":"created_at","KeyType":"RANGE"}
        ],
        "Projection": {"ProjectionType":"ALL"}
      },
      {
        "IndexName": "showtime-bookings-index",
        "KeySchema": [
          {"AttributeName":"showtime_id","KeyType":"HASH"},
          {"AttributeName":"created_at","KeyType":"RANGE"}
        ],
        "Projection": {"ProjectionType":"ALL"}
      },
      {
        "IndexName": "status-bookings-index",
        "KeySchema": [
          {"AttributeName":"status","KeyType":"HASH"},
          {"AttributeName":"created_at","KeyType":"RANGE"}
        ],
        "Projection": {"ProjectionType":"ALL"}
      }
  ]'

echo "Creating BOOKED_SEATS table..."

aws dynamodb create-table \
  $ENDPOINT $REGION \
  --table-name booked_seats \
  --attribute-definitions \
      AttributeName=id,AttributeType=S \
      AttributeName=booking_id,AttributeType=S \
      AttributeName=showtime_id,AttributeType=S \
      AttributeName=seat_id,AttributeType=S \
      AttributeName=status,AttributeType=S \
  --key-schema \
      AttributeName=id,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --global-secondary-indexes '[
      {
        "IndexName": "booking-seats-index",
        "KeySchema": [
          {"AttributeName":"booking_id","KeyType":"HASH"},
          {"AttributeName":"seat_id","KeyType":"RANGE"}
        ],
        "Projection": {"ProjectionType":"ALL"}
      },
      {
        "IndexName": "showtime-seats-index",
        "KeySchema": [
          {"AttributeName":"showtime_id","KeyType":"HASH"},
          {"AttributeName":"seat_id","KeyType":"RANGE"}
        ],
        "Projection": {"ProjectionType":"ALL"}
      },
      {
        "IndexName": "status-seats-index",
        "KeySchema": [
          {"AttributeName":"status","KeyType":"HASH"},
          {"AttributeName":"seat_id","KeyType":"RANGE"}
        ],
        "Projection": {"ProjectionType":"ALL"}
      }
  ]'

echo "✅ Tables created successfully"