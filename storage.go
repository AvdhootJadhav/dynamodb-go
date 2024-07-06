package main

import (
	"context"
	"log"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBStore struct {
	TableName string
	Client    *dynamodb.Client
}

type Storage interface {
	InsertAnime(Anime) error
	GetAnime(string) (Anime, error)
	DeleteAnime(string) error
}

func InitStore() (*DynamoDBStore, error) {
	endpoint := string("http://localhost:8000")
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("localhost"),
		config.WithSharedConfigProfile("default"),
	)

	cfg.BaseEndpoint = &endpoint

	if err != nil {
		return nil, err
	}
	store := DynamoDBStore{
		Client:    dynamodb.NewFromConfig(cfg),
		TableName: "anime",
	}
	return &store, nil
}

func (store DynamoDBStore) CheckTableExists(table string) bool {
	var tables []string
	tablePaginator := dynamodb.NewListTablesPaginator(store.Client, &dynamodb.ListTablesInput{})
	for tablePaginator.HasMorePages() {
		output, err := tablePaginator.NextPage(context.TODO())
		if err != nil {
			return false
		} else {
			tables = append(tables, output.TableNames...)
		}
	}
	return slices.Contains(tables, table)
}

func (store *DynamoDBStore) CreateTable() (*types.TableDescription, error) {
	tableOutput, err := store.Client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName: &store.TableName,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})

	if err != nil {
		return nil, err
	}
	waiter := dynamodb.NewTableExistsWaiter(store.Client)
	err = waiter.Wait(context.Background(), &dynamodb.DescribeTableInput{
		TableName: aws.String(store.TableName),
	}, 5*time.Minute)

	if err != nil {
		log.Fatalln("Wait for table exists failed due to", err)
	}

	return tableOutput.TableDescription, nil
}

func (store DynamoDBStore) InsertAnime(anime Anime) error {
	item, err := attributevalue.MarshalMap(anime)

	if err != nil {
		return err
	}

	_, err = store.Client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(store.TableName),
		Item:      item,
	})
	return err
}

func (store DynamoDBStore) GetAnime(id string) (Anime, error) {
	anime := new(Anime)
	temp, err := attributevalue.Marshal(id)
	if err != nil {
		return *anime, err
	}
	newId := map[string]types.AttributeValue{"id": temp}
	response, err := store.Client.GetItem(context.Background(), &dynamodb.GetItemInput{
		Key: newId, TableName: aws.String(store.TableName),
	})

	if err != nil {
		return *anime, err
	}

	err = attributevalue.UnmarshalMap(response.Item, &anime)

	if err != nil {
		return *anime, err
	}

	return *anime, nil
}

func (store DynamoDBStore) DeleteAnime(id string) error {
	temp, err := attributevalue.Marshal(id)

	if err != nil {
		return err
	}

	newId := map[string]types.AttributeValue{"id": temp}

	_, err = store.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: &store.TableName,
		Key:       newId,
	})
	return err
}
