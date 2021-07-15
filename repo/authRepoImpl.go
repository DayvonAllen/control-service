package repo

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	//"example.com/app/util"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepoImpl struct {
}

func(a AuthRepoImpl) Login(username string, password string, ip string, ips []string) (*domain.Admin, string, error) {
	var login domain.Authentication
	var admin domain.Admin

	conn := database.MongoConnectionPool.Get().(*database.Connection)
	database.MongoConnectionPool.Put(conn)

	opts := options.FindOne()
	err := conn.AdminCollection.FindOne(context.TODO(), bson.D{{"username",
		username}},opts).Decode(&admin)

	fmt.Println(admin)
	fmt.Println(password)

	if err != nil {
		return nil, "", fmt.Errorf("error finding by username")
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))

	if err != nil {
		return nil, "", fmt.Errorf("error comparing password")
	}

	token, err := login.GenerateJWT(admin)

	if err != nil {
		return nil, "", fmt.Errorf("error generating token")
	}

	go func() {
		filter := bson.D{{"username", username}}
		update := bson.D{{"$set", bson.D{{"lastLoginIp", ip}, {"lastLoginIps", ips}}}}

		_, err := conn.AdminCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			panic(err)
		}
		return
	}()

	return &admin, token, nil
}

func NewAuthRepoImpl() AuthRepoImpl {
	var authRepoImpl AuthRepoImpl

	return authRepoImpl
}