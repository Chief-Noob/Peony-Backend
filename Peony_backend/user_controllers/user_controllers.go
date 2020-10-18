package user_controllers

import (
	"Peony/Peony_backend/models/db"
	"Peony/Peony_backend/models/entity"
	"Peony/config"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

type GoogleRep struct {
	Id             string `json:"id"`
	Email          string `json:"email"`
	Verified_email bool   `json:"verified_email"`
	Picture        string `json:"picture"`
}

type CreateBodyForm struct {
	Student_number string `json:"student_number"`
	School         string `json:"school"`
	Email          string `json:"email"`
}

var cred = Credentials{
	Cid:     config.GetCredId(),
	Csecret: config.GetCredSct(),
}

var oauth2_config *oauth2.Config = &oauth2.Config{
	ClientID:     cred.Cid,
	ClientSecret: cred.Csecret,
	RedirectURL:  "http://localhost:8080/user/redir",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
	},
	Endpoint: google.Endpoint,
}

var jwtSecret = []byte(config.GetSecretKey())

func UserDetail(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "Here",
	})
}

func CreateUser(c *gin.Context) {
	client := db.GetConnection()
	collection := client.Database("Kebiao").Collection("user")

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Request body missing.",
		})
	}
	var body_form CreateBodyForm
	json.Unmarshal(body, &body_form)

	filter := bson.M{
		"email": body_form.Email,
	}
	if err != nil {
		log.Fatal(err)
	}

	var exist_user entity.User
	err = collection.FindOne(context.TODO(), filter).Decode(&exist_user)
	if err != nil {
		new_user := entity.User{
			body_form.Student_number,
			body_form.School,
			body_form.Email,
			[]bson.ObjectId{},
		}
		_, err := collection.InsertOne(context.TODO(), new_user)
		if err != nil {
			log.Fatal(err)
		}
		token := GetToken(&new_user)

		c.JSON(201, gin.H{
			"message": "New user created!",
			"token":   token,
		})
		return
	}

	c.JSON(409, gin.H{
		"message": "User already exist.",
	})
	return
}

func UserGmail(c *gin.Context) {
	code := c.Request.URL.Query().Get("code")
	google_access_tok, err := oauth2_config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
		return
	}

	var client = oauth2_config.Client(oauth2.NoContext, google_access_tok)
	response, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": "Invalid google_access_tok."})
	}

	res_byte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Json unmarshal failed.")
	}

	var res_struct GoogleRep
	json.Unmarshal(res_byte, &res_struct)

	db_client := db.GetConnection()
	collection := db_client.Database("Kebiao").Collection("user")

	filter := bson.M{
		"email": res_struct.Email,
	}
	if err != nil {
		log.Fatal(err)
	}

	var exist_user entity.User
	err = collection.FindOne(context.TODO(), filter).Decode(&exist_user)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "Please created new user.",
		})
		return
	}
	token := GetToken(&exist_user)

	c.JSON(200, gin.H{
		"token": token,
	})
	return
}

func GetToken(u *(entity.User)) string {
	now := time.Now()
	jwtId := (*u).Student_number + strconv.FormatInt(now.Unix(), 10)
	claims := entity.Claims{
		(*u).Student_number,
		(*u).School,
		(*u).Email,
		jwt.StandardClaims{
			Audience:  (*u).Student_number,
			ExpiresAt: now.Add(86400 * time.Second).Unix(),
			Id:        jwtId,
			IssuedAt:  now.Unix(),
			Issuer:    "ginJWT",
			NotBefore: now.Add(10 * time.Second).Unix(),
			Subject:   (*u).Student_number,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		fmt.Println("here pro")
		log.Fatal(err)
	}
	return token
}

func RandToken(l int) (string, error) {
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func AuthHandler(c *gin.Context) {
	state, err := RandToken(32)
	if err != nil {
		log.Fatal(err)
		return
	}
	red_url := oauth2_config.AuthCodeURL(state)
	c.JSON(200, gin.H{
		"redirect_url": red_url,
	})
	return
}
