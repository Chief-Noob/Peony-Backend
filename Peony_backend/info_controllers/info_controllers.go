package info_controllers

import (
	"Peony/Peony_backend/deserializers"
	_ "Peony/Peony_backend/deserializers"
	"Peony/Peony_backend/models/db"
	"Peony/Peony_backend/models/entity"
	"Peony/config"
	"context"
	"io/ioutil"
	"strings"
	_ "time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
)

var jwt_secret = []byte(config.GetSecretKey())

func CreateInfo(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(404, gin.H{
			"error": "NO REQUEST BODY.",
		})
		return
	}

	var body_info entity.Info

	body_info, err = deserializers.InfoDeserializer(body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
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
	if err == nil {
		c.JSON(409, gin.H{
			"error": "INFO ALREADY EXIST.",
		})
		return
	}

	insert_result, err := collection.InsertOne(context.TODO(), body_info)
	if err != nil {
		c.JSON(405, gin.H{
			"error": "DB INSERT FAIL.",
		})
		return
	}

	auth_header := c.Request.Header.Get("Authetication")
	token := strings.Split(auth_header, " ")[1]
	token_claims, err := jwt.ParseWithClaims(token, &entity.Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return jwt_secret, nil
	})
	claims, _ := token_claims.Claims.(*entity.Claims)
	user_filter := bson.M{
		"email": claims.Email,
	}

	user_update := bson.D{
		{"$addToSet", bson.M{"infolist": insert_result.InsertedID}},
	}
	collection = client.Database("Kebiao").Collection("user")
	_, err = collection.UpdateOne(
		context.TODO(),
		user_filter,
		user_update,
	)
	if err != nil {
		c.JSON(406, gin.H{
			"error": "USER UPDATE FAIL.",
		})
		return
	}

	c.JSON(201, body_info)
	return
}

func InfoDetail(c *gin.Context) {
	info_id := c.DefaultQuery("info_id", "None")
	var info entity.Info

	if info_id != "None" {
		client := db.GetConnection()
		collection := client.Database("Kebiao").Collection("info")
		info_id_pure, _ := primitive.ObjectIDFromHex(info_id)
		filter := bson.M{
			"_id": info_id_pure,
		}
		err := collection.FindOne(context.TODO(), filter).Decode(&info)
		if err != nil {
			c.JSON(409, gin.H{
				"error": "INFO NOT EXIST.",
			})
			return
		}
		c.JSON(200, info)
		return
	}

	course_number := c.DefaultQuery("course_number", "None")
	school := c.DefaultQuery("school", "None")
	if course_number != "None" && school != "None" {
		client := db.GetConnection()
		collection := client.Database("Kebiao").Collection("info")
		filter := bson.M{
			"coursenumber": course_number,
			"school":       school,
		}
		var all_info []bson.M
		cursor, err := collection.Find(context.TODO(), filter)
		if err != nil {
			c.JSON(409, gin.H{
				"error": "INFO NOT EXISTS.",
			})
			return
		}
		err = cursor.All(context.TODO(), &all_info)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "CURSOR ERROR.",
			})
			return
		}
		c.JSON(200, all_info)
		return
	}
	c.JSON(400, gin.H{
		"error": "PARAMS MALFORMED.",
	})
	return
}

func UpdateInfo(c *gin.Context) {
	info_id_raw := c.Param("info_id")
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "REQUEST BODY MALFORMED.",
		})
		return
	}
	var body_info entity.Info
	body_info, err = deserializers.InfoDeserializer(body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	info_id, err := primitive.ObjectIDFromHex(info_id_raw)
	if err != nil {
		c.JSON(406, gin.H{
			"error": "INFO_ID MALFORMED.",
		})
		return
	}

	client := db.GetConnection()
	collection := client.Database("Kebiao").Collection("info")
	filter := bson.M{
		"_id": info_id,
	}
	update := bson.D{{"$set", bson.M{
		"coursenumber": body_info.CourseNumber,
		"school":       body_info.School,
		"fieldtitle":   body_info.FieldTitle,
		"fieldcontent": body_info.FieldContent,
		"origin":       body_info.Origin,
		"starttime":    body_info.StartTime,
		"endtime":      body_info.EndTime,
	}},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(404, gin.H{
			"error": "USER UPDATE FAIL.",
		})
		return
	}

	c.JSON(200, body_info)
	return
}
