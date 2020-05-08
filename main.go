package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vinhut/posted/helpers"
	"github.com/vinhut/posted/models"
	"github.com/vinhut/posted/services"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"encoding/json"
	"fmt"
	"time"
)

var SERVICE_NAME = "post-service"

type UserAuthData struct {
	Uid     string
	Email   string
	Role    string
	Created string
}

func checkUser(authservice services.AuthService, token string) (*UserAuthData, error) {

	data := &UserAuthData{}
	user_data, auth_error := authservice.Check(SERVICE_NAME, token)
	if auth_error != nil {
		return data, auth_error
	}

	if err := json.Unmarshal([]byte(user_data), data); err != nil {
		fmt.Println(err)
		return data, err
	}

	return data, nil

}

func setupRouter(postdb models.PostDatabase, authservice services.AuthService) *gin.Engine {

	router := gin.Default()

	router.GET(SERVICE_NAME+"/ping", func(c *gin.Context) {
		c.String(200, "OK")
	})

	router.GET(SERVICE_NAME+"/post", func(c *gin.Context) {

		value, err := c.Cookie("token")
		post_id, _ := c.GetQuery("postid")
		if err != nil {
			panic("failed get token")
		}
		_, check_err := checkUser(authservice, value)
		if check_err != nil {
			panic("error check user")
		}

		result := &models.Post{}
		find_err := postdb.Find("_id", post_id, result)
		if find_err != nil {
			panic("can't find post")
		}

		post_json, json_err := json.Marshal(result)
		if json_err != nil {
			panic("marshal json fail")
		}

		c.String(200, string(post_json))

	})

	router.POST(SERVICE_NAME+"/post", func(c *gin.Context) {

		value, err := c.Cookie("token")
		if err != nil {
			panic("failed get token")
		}
		user_data, check_err := checkUser(authservice, value)
		if check_err != nil {
			panic("error check user")
		}
		img_url := c.PostForm("img_url")
		post_caption := c.PostForm("post_caption")

		new_post := &models.Post{

			Postid:       primitive.NewObjectIDFromTimestamp(time.Now()),
			Uid:          user_data.Uid,
			Imageurl:     img_url,
			Caption:      post_caption,
			Likecount:    0,
			Private:      false,
			Commentcount: 0,
			Viewcount:    0,
			Created:      time.Now(),
			Tag:          make([]string, 1),
		}

		postdb.Create(new_post)

	})

	router.DELETE(SERVICE_NAME+"/post", func(c *gin.Context) {

		value, err := c.Cookie("token")
		post_id, _ := c.GetQuery("postid")
		if err != nil {
			panic("failed get token")
		}
		_, check_err := checkUser(authservice, value)
		if check_err != nil {
			panic("error check user")
		}

		_, delete_err := postdb.Delete(post_id)
		if delete_err != nil {
			panic("can't delete post")
		}

		c.String(200, "deleted")

	})

	router.GET(SERVICE_NAME+"/allpost", func(c *gin.Context) {

		result, findall_err := postdb.FindAll()
		if findall_err != nil {
			panic("error getting all post")
		}

		allid_json, json_err := json.Marshal(result)
		if json_err != nil {
			panic("marshal json fail")
		}
		c.String(200, `{ "results": `+string(allid_json)+`}`)

	})

	return router

}

func main() {

	mongo_layer := helpers.NewMongoDatabase()
	postdb := models.NewPostDatabase(mongo_layer)
	authservice := services.NewUserAuthService()
	router := setupRouter(postdb, authservice)
	router.Run(":8080")

}
