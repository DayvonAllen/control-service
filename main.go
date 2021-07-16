package main

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	"example.com/app/event-consumer"
	"example.com/app/router"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"os/signal"
)

func init() {
	// create database connection instance for first time
	go event_consumer.KafkaConsumerGroup()
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	adminSearch := new(domain.Admin)
	err := conn.AdminCollection.FindOne(context.TODO(), bson.M{"username": "admin"}).Decode(adminSearch)

	if err != nil  {
		if err == mongo.ErrNoDocuments {
			admin := domain.Admin{Username: "admin", Password: "password"}
			admin.Id = primitive.NewObjectID()
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
			admin.Password = string(hashedPassword)

			_, err := conn.AdminCollection.InsertOne(context.TODO(), &admin)

			if err != nil {
				panic("error processing data")
			}
			return
		}
		panic(err)
	}
}

func main() {
	app := router.Setup()

	// graceful shutdown on signal interrupts
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		_ = <-c
		fmt.Println("Shutting down...")
		_ = app.Shutdown()
	}()

	if err := app.Listen(":8084"); err != nil {
		log.Panic(err)
	}
}