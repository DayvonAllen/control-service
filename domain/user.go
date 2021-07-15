package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id                          primitive.ObjectID   `bson:"_id" json:"id"`
	Username                    string               `bson:"username" json:"username"`
	Email                       string               `bson:"email" json:"email"`
	Password                    string               `bson:"password" json:"-"`
	CurrentTagLine              string               `bson:"currentTagLine" json:"CurrentTagLine"`
	ProfilePictureUrl           string               `bson:"profilePictureUrl" json:"profilePictureUrl"`
	ProfileBackgroundPictureUrl string               `bson:"profileBackgroundPictureUrl" json:"profileBackgroundPictureUrl"`
	CurrentBadgeUrl             string               `bson:"currentBadgeUrl" json:"currentBadgeUrl"`
	UnlockedBadgesUrls          []string             `bson:"unlockedBadgesUrls" json:"unlockedBadgesUrls"`
	BlockList                   []string `bson:"blockList" json:"blockList"`
	BlockByList                 []string `bson:"blockByList" json:"blockByList"`
	FlagCount                   []primitive.ObjectID `bson:"flagCount" json:"-"`
	Followers                   []string             `bson:"followers" json:"followers"`
	Following                   []string             `bson:"following" json:"following"`
	FollowerCount               int                  `bson:"followerCount" json:"followerCount"`
	DisplayFollowerCount        bool                 `bson:"displayFollowerCount" json:"displayFollowerCount"`
	ProfileIsViewable           bool                 `bson:"profileIsViewable" json:"profileIsViewable"`
	IsLocked                    bool                 `bson:"isLocked" json:"-"`
	IsVerified                  bool                 `bson:"isVerified" json:"isVerified"`
	AcceptMessages              bool                 `bson:"acceptMessages" json:"acceptMessages"`
	LastLoginIp					string				 `bson:"lastLoginIp" json:"-"`
	LastLoginIps				[]string			 `bson:"lastLoginIps" json:"-"`
	CreatedAt                   time.Time            `bson:"createdAt" json:"-"`
	UpdatedAt                   time.Time            `bson:"updatedAt" json:"-"`
}

type UserDto struct {
	Id                          primitive.ObjectID   `bson:"_id" json:"-"`
	Email                       string               `json:"email"`
	Username                    string               `json:"username"`
	CurrentTagLine              string               `json:"currentTagLine"`
	UnlockedTagLine             []string             `json:"unlockedTagLine"`
	ProfilePictureUrl           string               `json:"profilePictureUrl"`
	ProfileBackgroundPictureUrl string               `json:"profileBackgroundPictureUrl"`
	CurrentBadgeUrl             string               `json:"currentBadgeUrl"`
	ProfileIsViewable           bool                 `json:"profileIsViewable"`
	AcceptMessages              bool                 `json:"acceptMessages"`
	FollowerCount               int                  `json:"followerCount"`
	DisplayFollowerCount        bool                 `json:"displayFollowerCount"`
	Followers                   []string             `bson:"followers" json:"-"`
	Following                   []string             `bson:"following" json:"-"`
}

type UserResponse struct {
	Users       *[]UserDto
	CurrentPage string
}