package info_controllers

import (
	"Peony/Peony_backend/models/db"
	"Peony/Peony_backend/models/entity"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateInfo(c *gin.Context) {
	body_json, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(404, gin.H{
			"error": "NO REQUEST BODY.",
		})
		return
	}
	var body_info entity.Info
	json.Unmarshal(body_json, &body_info)
	fmt.Println(body_info.CourseNumber)

	client := db.GetConnection()
	collection := client.Database("Kebiao").Collection("info")
	filter := bson.M{
		"coursenumber": body_info.CourseNumber,
		"school":       body_info.School,
		"fieldtitle":   body_info.FieldTitle,
		"fieldcontent": body_info.FieldContent,
		"origin":       body_info.Origin,
		"starttime":    body_info.StartTime,
		"endtime":      body_info.EndTime,
	}
	var exist_info entity.Info
	err = collection.FindOne(context.TODO(), filter).Decode(&exist_info)
	if err != nil {
		new_info := entity.Info{
			body_info.CourseNumber,
			body_info.School,
			body_info.FieldTitle,
			body_info.FieldContent,
			body_info.Origin,
			body_info.StartTime,
			body_info.EndTime,
		}
		_, err = collection.InsertOne(context.TODO(), new_info)
		if err != nil {
			c.JSON(405, gin.H{
				"error": "DB INSERT FAIL.",
			})
			return
		}

		c.JSON(201, gin.H{
			"coursenumber": new_info.CourseNumber,
			"school":       new_info.School,
			"fieldtitle":   new_info.FieldTitle,
			"fieldcontent": new_info.FieldContent,
			"origin":       new_info.Origin,
			"starttime":    new_info.StartTime,
			"endtime":      new_info.EndTime,
		})
		return
	}
	c.JSON(409, gin.H{
		"error": "INFO ALREADY EXIST.",
	})
	return
}
