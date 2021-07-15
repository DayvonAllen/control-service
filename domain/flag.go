package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// Flag todo validate struct
type Flag struct {
	Id              primitive.ObjectID `bson:"_id" json:"-"`
	FlaggerID       primitive.ObjectID `bson:"flaggerID" json:"-"`
	FlaggedResource primitive.ObjectID `bson:"flaggedResource" json:"-"`
	Reason          string             `bson:"reason" json:"reason"`
}
