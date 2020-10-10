package user_controllers

import (
	"Peony/Peony_backend/models/db"
	"Peony/Peony_backend/models/entity"
	"Peony/config"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

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

type Claims struct {
	Student_number string `json:"student_number"`
	School         string `json:"school"`
	jwt.StandardClaims
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

func CreateUser(c *gin.Context) {
	client := db.GetConnection()
	collection := client.Database("Kebiao").Collection("user")

	new_user := entity.User{
		"0812253",
		"nctu",
		"aaaa102234@gmail.com",
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

	c.JSON(http.StatusOK, gin.H{
		"email": res_struct.Email,
	})
	return
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
}
