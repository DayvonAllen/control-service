package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Reply struct {
	Id                  primitive.ObjectID `bson:"_id" json:"id"`
	ResourceId 			primitive.ObjectID `bson:"resourceId" json:"-"`
	Content             string             `bson:"content" json:"content"`
	AuthorUsername      string             `bson:"authorUsername" json:"authorUsername"`
	Likes          		[]string           `bson:"likes" json:"-"`
	Dislikes       		[]string           `bson:"dislikes" json:"-"`
	LikeCount           int                `bson:"likeCount" json:"likeCount"`
	DislikeCount        int                `bson:"dislikeCount" json:"dislikeCount"`
	Edited              bool               `bson:"edited" json:"edited"`
	CurrentUserLiked    bool               `bson:"currentUserLiked" json:"currentUserLiked"`
	CurrentUserDisLiked bool               `bson:"currentUserDisLiked" json:"currentUserDisLiked"`
	CreatedAt           time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt           time.Time          `bson:"updatedAt" json:"updatedAt"`
	CreatedDate         string             `bson:"createdDate" json:"createdDate"`
	UpdatedDate         string             `bson:"updatedDate" json:"updatedDate"`
}
