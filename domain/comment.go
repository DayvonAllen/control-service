package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Comment todo validate struct
type Comment struct {
	Id             primitive.ObjectID `bson:"_id" json:"-"`
	ResourceId     primitive.ObjectID `bson:"resourceId" json:"-"`
	Content        string             `bson:"content" json:"content"`
	AuthorUsername string             `bson:"authorUsername" json:"-"`
	Edited         bool               `bson:"edited" json:"-"`
	Likes          []string           `bson:"likes" json:"-"`
	Dislikes       []string           `bson:"dislikes" json:"-"`
	LikeCount      int                `bson:"likeCount" json:"-"`
	DislikeCount   int                `bson:"dislikeCount" json:"-"`
	CreatedAt      time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"-"`
	CreatedDate    string             `bson:"createdDate" json:"-"`
	UpdatedDate    string             `bson:"updatedDate" json:"-"`
}

type CommentDto struct {
	Id                  primitive.ObjectID `bson:"_id" json:"id"`
	Content             string             `bson:"content" json:"content"`
	AuthorUsername      string             `bson:"authorUsername" json:"authorUsername"`
	Likes          []string           `bson:"likes" json:"-"`
	Dislikes       []string           `bson:"dislikes" json:"-"`
	LikeCount           int                `bson:"likeCount" json:"likeCount"`
	DislikeCount        int                `bson:"dislikeCount" json:"dislikeCount"`
	Edited              bool               `bson:"edited" json:"edited"`
	Replies             *[]Reply      `bson:"replies" json:"replies"`
	CurrentUserLiked    bool               `bson:"currentUserLiked" json:"currentUserLiked"`
	CurrentUserDisLiked bool               `bson:"currentUserDisLiked" json:"currentUserDisLiked"`
	CreatedAt           time.Time          `json:"createdAt"`
	UpdatedAt           time.Time          `json:"updatedAt"`
	CreatedDate         string             `json:"createdDate"`
	UpdatedDate         string             `json:"updatedDate"`
}

