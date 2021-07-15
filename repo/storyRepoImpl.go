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
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"strconv"
	"sync"
	"time"
)

type StoryRepoImpl struct {
	Story             domain.Story
	StoryDto          domain.StoryDto
	StoryList         []domain.Story
	StoryDtoList      []domain.StoryDto
}

func (s StoryRepoImpl) FindAll(page string, newStoriesQuery bool) (*[]domain.Story, error) {
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

	if newStoriesQuery {
		findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	cur, err := conn.StoryCollection.Find(context.TODO(), bson.M{}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &s.StoryList); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	return &s.StoryList, nil
}

func (s StoryRepoImpl) FindById(storyID primitive.ObjectID) (*domain.StoryDto, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	err := conn.StoryCollection.FindOne(context.TODO(), bson.D{{"_id", storyID}}).Decode(&s.StoryDto)

	if err != nil {
		return nil, err
	}
	return &s.StoryDto, nil
}

func (s StoryRepoImpl) Create(story *domain.Story) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	story.Id = primitive.NewObjectID()

	_, err := conn.StoryCollection.InsertOne(context.TODO(), &story)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (s StoryRepoImpl) UpdateById(id primitive.ObjectID, newContent string, newTitle string, username string, tags *[]domain.Tag, updated bool) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	filter := bson.D{{"_id", id}, {"authorUsername", username}}
	update := bson.D{{"$set",
		bson.D{{"content", newContent},
			{"title", newTitle},
			{"updatedAt", time.Now()},
			{"tags", tags},
			{"updated", updated},
		},
	}}

	_, err := conn.StoryCollection.UpdateOne(context.TODO(),
		filter, update)

	if err != nil {
		return fmt.Errorf("you can't update a story you didn't write")
	}

	return nil
}


func (s StoryRepoImpl) DeleteById(id primitive.ObjectID) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)
	// sets mongo's read and write concerns
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	// set up for a transaction
	session, err := conn.StartSession()

	if err != nil {
		panic(err)
	}

	defer session.EndSession(context.Background())

	// execute this code in a logical transaction
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		var wg sync.WaitGroup
		wg.Add(3)

		go func() {
			defer wg.Done()
			res, err := conn.StoryCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}})

			if err != nil {
				panic(err)
			}

			if res.DeletedCount == 0 {
				panic(fmt.Errorf("failed to delete story"))
			}
			return
		}()

		go func() {
			defer wg.Done()
			_, err = conn.FlagCollection.DeleteMany(context.TODO(), bson.D{{"flaggedResource", id}})

			if err != nil {
				panic(err)
			}
			return
		}()

		go func() {
			defer wg.Done()
			err = CommentRepoImpl{}.DeleteManyById(id)

			if err != nil {
				panic(err)
			}
			return
		}()

		wg.Wait()

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return err
	}

	return nil
}

func NewStoryRepoImpl() StoryRepoImpl {
	var storyRepoImpl StoryRepoImpl

	return storyRepoImpl
}
