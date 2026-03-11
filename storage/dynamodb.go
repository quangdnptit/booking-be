package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"booking-be/models"
)

// DynamoDB implements Store interface with DynamoDB
type DynamoDB struct {
	client        *dynamodb.Client
	roomsTable    string
	bookingsTable string
}

// NewDynamoDBStore creates and initializes a DynamoDB store
func NewDynamoDBStore(ctx context.Context, endpoint string) (*DynamoDB, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	// Allow custom endpoint for local development
	clientOptions := []func(*dynamodb.Options){}
	if endpoint != "" {
		clientOptions = append(clientOptions, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}

	client := dynamodb.NewFromConfig(cfg, clientOptions...)

	store := &DynamoDB{
		client:        client,
		roomsTable:    "rooms",
		bookingsTable: "bookings",
	}

	// Initialize tables if they don't exist
	if err := store.initializeTables(ctx); err != nil {
		return nil, err
	}

	// Seed initial rooms data
	if err := store.seedRooms(ctx); err != nil {
		return nil, err
	}

	return store, nil
}

// initializeTables creates the required tables
func (s *DynamoDB) initializeTables(ctx context.Context) error {
	tables := []struct {
		name      string
		partition string
	}{
		{s.roomsTable, "id"},
		{s.bookingsTable, "id"},
	}

	for _, t := range tables {
		if err := s.createTableIfNotExists(ctx, t.name, t.partition); err != nil {
			return err
		}
	}

	return nil
}

// createTableIfNotExists creates a table if it doesn't already exist
func (s *DynamoDB) createTableIfNotExists(ctx context.Context, tableName, partitionKey string) error {
	_, err := s.client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})

	if err == nil {
		return nil // Table already exists
	}

	// Create table
	_, err = s.client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String(partitionKey),
				KeyType:       types.KeyTypeHash,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String(partitionKey),
				AttributeType: types.ScalarAttributeTypeS, // String
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})

	if err != nil {
		return fmt.Errorf("unable to create table %s: %w", tableName, err)
	}

	return nil
}

// seedRooms seeds initial room data
func (s *DynamoDB) seedRooms(ctx context.Context) error {
	rooms := []models.Room{
		{ID: "room-1", Name: "Deluxe Suite", Capacity: 2, PricePerNight: 150.00},
		{ID: "room-2", Name: "Family Room", Capacity: 4, PricePerNight: 200.00},
	}

	for _, room := range rooms {
		item, err := attributevalue.MarshalMap(room)
		if err != nil {
			return fmt.Errorf("failed to marshal room: %w", err)
		}

		_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(s.roomsTable),
			Item:      item,
		})

		if err != nil {
			return fmt.Errorf("failed to put room: %w", err)
		}
	}

	return nil
}

// GetRooms returns all rooms from DynamoDB
func (s *DynamoDB) GetRooms() []models.Room {
	ctx := context.Background()
	result, err := s.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(s.roomsTable),
	})

	if err != nil {
		return []models.Room{}
	}

	var rooms []models.Room
	err = attributevalue.UnmarshalListOfMaps(result.Items, &rooms)
	if err != nil {
		return []models.Room{}
	}

	return rooms
}

// GetRoomByID returns a room by its ID from DynamoDB
func (s *DynamoDB) GetRoomByID(id string) (*models.Room, bool) {
	ctx := context.Background()
	result, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.roomsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil || result.Item == nil {
		return nil, false
	}

	var room models.Room
	err = attributevalue.UnmarshalMap(result.Item, &room)
	if err != nil {
		return nil, false
	}

	return &room, true
}

// GetBookings returns all bookings from DynamoDB
func (s *DynamoDB) GetBookings() []models.Booking {
	ctx := context.Background()
	result, err := s.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(s.bookingsTable),
	})

	if err != nil {
		return []models.Booking{}
	}

	var bookings []models.Booking
	err = attributevalue.UnmarshalListOfMaps(result.Items, &bookings)
	if err != nil {
		return []models.Booking{}
	}

	return bookings
}

// GetBookingByID returns a booking by its ID from DynamoDB
func (s *DynamoDB) GetBookingByID(id string) (*models.Booking, bool) {
	ctx := context.Background()
	result, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.bookingsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil || result.Item == nil {
		return nil, false
	}

	var booking models.Booking
	err = attributevalue.UnmarshalMap(result.Item, &booking)
	if err != nil {
		return nil, false
	}

	return &booking, true
}

// CreateBooking creates and stores a new booking in DynamoDB
func (s *DynamoDB) CreateBooking(booking models.Booking) models.Booking {
	ctx := context.Background()
	item, err := attributevalue.MarshalMap(booking)
	if err != nil {
		return booking
	}

	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.bookingsTable),
		Item:      item,
	})

	if err != nil {
		return booking
	}

	return booking
}
