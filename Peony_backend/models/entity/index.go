package entity

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	StudentNumber string               `json:”studentnumber”`
	School        string               `json:”school”`
	Email         string               `json:”email”`
	InfoList      []primitive.ObjectID `json:”infolist”`
}

type UserWithId struct {
	StudentNumber string               `json:”studentnumber”`
	School        string               `json:”school”`
	Email         string               `json:”email”`
	InfoList      []primitive.ObjectID `json:”infolist”`
	Id            primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
}

type Info struct {
	CourseNumber string    `json:”coursenumber”`
	School       string    `json:”school”`
	FieldTitle   string    `json:”fieldtitle"`
	FieldContent string    `json:”fieldcontent”`
	Origin       string    `json:”origin”`
	StartTime    time.Time `json:”starttime”`
	EndTime      time.Time `json:”endtime”`
}

type Claims struct {
	StudentNumber string `json:"studentnumber"`
	School        string `json:"school"`
	Email         string `json:"email"`
	jwt.StandardClaims
}
