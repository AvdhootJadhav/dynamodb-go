package main

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	store, err := InitStore()

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connected to DB!!")
	res, err := store.CreateTable()

	if err != nil {
		log.Fatalln(err)
	}

	if res.TableStatus == types.TableStatusActive {
		log.Println("Table created!!")
	}
	server := NewServer(":8080", store)
	server.Run()
}
