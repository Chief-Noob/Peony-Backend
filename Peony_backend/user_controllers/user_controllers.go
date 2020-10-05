package user_controllers

import (
	"Peony/Peony_backend/models/db"
	"Peony/Peony_backend/models/entity"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func CreateUser(c *gin.Context) {
	client := db.GetConnection()
	collection := client.Database("Kebiao").Collection("user")

	new_user := entity.User{
		"0812253",
		"nctu",
		[]bson.ObjectId{
			bson.NewObjectId(),
			bson.NewObjectId(),
		},
	}

	insertResult, err := collection.InsertOne(context.TODO(), new_user)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted a single document: ", insertResult.InsertedID)
}
