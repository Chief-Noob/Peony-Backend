package entity

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
)

type User struct {
	Student_number string          `json:”student_number”`
	School         string          `json:”school”`
	Email          string          `json:”email”`
	Info_list      []bson.ObjectId `json:”info_list”`
}

type Info struct {
	Course_number string     `json:”course_number”`
	School        string     `json:”school”`
	Field_title   string     `json:”field_title"`
	Field_content string     `json:”field_content”`
	Origin        string     `json:”origin”`
	Time          *time.Time `json:”time”`
}

type Claims struct {
	Student_number string `json:"student_number"`
	School         string `json:"school"`
	Email          string `json:"email"`
	jwt.StandardClaims
}
