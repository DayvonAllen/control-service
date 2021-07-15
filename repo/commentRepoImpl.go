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

type CommentRepoImpl struct {
	Comment        domain.Comment
	CommentDto     domain.CommentDto
	Reply          domain.Reply
	CommentList    []domain.Comment
	CommentDtoList []domain.CommentDto
}

func (c CommentRepoImpl) Create(comment *domain.Comment) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	story := new(domain.Story)

	err := conn.StoryCollection.FindOne(context.TODO(), bson.D{{"_id", comment.ResourceId}}).Decode(&story)

	if err != nil {
		return fmt.Errorf("resource not found")
	}

	_, err = conn.CommentsCollection.InsertOne(context.TODO(), &comment)

	if err != nil {
		return err
	}

	go func() {
		event := new(domain.Event)
		event.Action = "comment on story"
		event.Target = comment.ResourceId.String()
		event.ResourceId = comment.ResourceId
		event.ActorUsername = comment.AuthorUsername
		event.Message = comment.AuthorUsername + " commented on a story with the ID:" + comment.ResourceId.String()
		err = SendEventMessage(event, 0)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}
	}()

	return nil
}

func (c CommentRepoImpl) UpdateById(id primitive.ObjectID, newContent string, edited bool, updatedTime time.Time, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}, {"authorUsername", username}}
	update := bson.D{{"$set", bson.D{{"content", newContent}, {"edited", edited},
		{"updatedTime", updatedTime}}}}

	err := conn.CommentsCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&c.Comment)

	if err != nil {
		return fmt.Errorf("cannot update comment that you didn't write")
	}

	return nil
}

func (c CommentRepoImpl) DeleteById(id primitive.ObjectID) error {
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

			res, err := conn.CommentsCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}})

			if err != nil {
				panic(err)
			}

			if res.DeletedCount == 0 {
				panic(fmt.Errorf("you can't delete a comment that you didn't create"))
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

			_, err = conn.RepliesCollection.DeleteMany(context.TODO(), bson.D{{"resourceId", id}})
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
		return fmt.Errorf("failed to delete reply")
	}

	return nil
}

func (c CommentRepoImpl) DeleteManyById(id primitive.ObjectID) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	err := conn.CommentsCollection.FindOne(context.TODO(), bson.D{{"resourceId", id}}).Decode(&c.Comment)

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

			res, err := conn.CommentsCollection.DeleteMany(context.TODO(), bson.D{{"resourceId", id}})
			if err != nil {
				panic(err)
			}

			if res.DeletedCount == 0 {
				panic(err)
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

			_, err = conn.RepliesCollection.DeleteMany(context.TODO(), bson.D{{"resourceId", c.Comment.Id}})
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

func NewCommentRepoImpl() CommentRepoImpl {
	var commentRepoImpl CommentRepoImpl

	return commentRepoImpl
}
