package repo

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
)

type UserRepoImpl struct {
	users        []domain.User
	user         domain.User
	userDto      domain.UserDto
	userDtoList  []domain.UserDto
	userResponse domain.UserResponse
}

func (u UserRepoImpl) FindAll(page string, ctx context.Context) (*domain.UserResponse, error) {

	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}
	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	// Get all users
	cur, err := conn.UserCollection.Find(ctx, bson.M{}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(ctx, &u.userDtoList); err != nil {
		log.Fatal(err)
	}

	u.userResponse = domain.UserResponse{Users: &u.userDtoList, CurrentPage: page}

	return &u.userResponse, nil
}

func (u UserRepoImpl) Create(user *domain.User) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	cur, err := conn.UserCollection.Find(context.TODO(), bson.M{
		"$or": []interface{}{
			bson.M{"email": user.Email},
			bson.M{"username": user.Username},
		},
	})

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	if !cur.Next(context.TODO()) {
		_, err = conn.UserCollection.InsertOne(context.TODO(), &user)

		if err != nil {
			return fmt.Errorf("error processing data")
		}

		return nil
	}

	return fmt.Errorf("user already exists")
}

func (u UserRepoImpl) UpdateByID(user *domain.User) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", user}}

	conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return nil
}

func (u UserRepoImpl) FindByUsername(username string) (*domain.UserDto, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	err := conn.UserCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&u.userDto)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("cannot find user")
		}
		return nil, fmt.Errorf("error processing data")
	}

	return &u.userDto, nil
}

func (u UserRepoImpl) DeleteByID(id primitive.ObjectID) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	_, err := conn.UserCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}})

	if err != nil {
		return err
	}

	return nil
}

func NewUserRepoImpl() UserRepoImpl {
	var userRepoImpl UserRepoImpl

	return userRepoImpl
}
