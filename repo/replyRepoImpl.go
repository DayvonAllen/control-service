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
	"sync"
	"time"
)

type ReplyRepoImpl struct {
	Reply        domain.Reply
	ReplyList    []domain.Reply
}

func (r ReplyRepoImpl) Create(comment *domain.Reply) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	commentObj := new(domain.Comment)

	err := conn.CommentsCollection.FindOne(context.TODO(), bson.D{{"_id", comment.ResourceId}}).Decode(&commentObj)

	if err != nil {
		return fmt.Errorf("resource not found")
	}

	_, err = conn.RepliesCollection.InsertOne(context.TODO(), &comment)

	if err != nil {
		return err
	}

	return nil
}

func (r ReplyRepoImpl) UpdateById(id primitive.ObjectID, newContent string, edited bool, updatedTime time.Time) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"content", newContent}, {"edited", edited},
		{"updatedTime", updatedTime}}}}

	err := conn.RepliesCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&r.Reply)

	if err != nil {
		return fmt.Errorf("cannot update comment that you didn't write")
	}

	return nil
}

func (r ReplyRepoImpl) DeleteById(id primitive.ObjectID, username string) error {
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
		wg.Add(2)

		go func() {
			defer wg.Done()

			res, err := conn.RepliesCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}, {"authorUsername", username}})

			if err != nil {
				panic(err)
			}

			if res.DeletedCount == 0  {
				panic(fmt.Errorf("failed to delete reply"))
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

		wg.Wait()

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return err
	}

	return nil
}

func NewReplyRepoImpl() ReplyRepoImpl {
	var replyRepoImpl ReplyRepoImpl

	return replyRepoImpl
}