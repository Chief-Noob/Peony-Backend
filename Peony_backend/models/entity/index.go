package entity

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type User struct {
	Student_number string          `json:”student_number,omitempty”`
	School         string          `json:”school,omitempty”`
	Info_list      []bson.ObjectId `json:”info_list,omitempty”`
}

type Info struct {
	Course_number string     `json:”course_number,omitempty”`
	School        string     `json:”school,omitempty”`
	Field_title   string     `json:”field_title,omitempty”`
	Field_content string     `json:”field_content,omitempty”`
	Origin        string     `json:”origin,omitempty”`
	Time          *time.Time `json:”time,omitempty”`
}
