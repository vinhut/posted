package main

import (
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"github.com/vinhut/posted/helpers"
	"github.com/vinhut/posted/models"
	"github.com/vinhut/posted/services"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var SERVICE_NAME = "post-service"

func checkUser(authservice services.AuthService, token string) (map[string]interface{}, error) {

	var data map[string]interface{}
	user_data, auth_error := authservice.Check(SERVICE_NAME, token)
	if auth_error != nil {
		return data, auth_error
	}

	if err := json.Unmarshal([]byte(user_data), &data); err != nil {
		fmt.Println(err)
		return data, err
	}

	return data, nil

}

func setupRouter(postdb models.PostDatabase, authservice services.AuthService) *gin.Engine {

	var JAEGER_COLLECTOR_ENDPOINT = os.Getenv("JAEGER_COLLECTOR_ENDPOINT")
	cfg := jaegercfg.Configuration{
		ServiceName: "post-service",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: JAEGER_COLLECTOR_ENDPOINT,
		},
	}
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory
	tracer, _, _ := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	opentracing.SetGlobalTracer(tracer)

	router := gin.Default()

	router.GET(SERVICE_NAME+"/ping", func(c *gin.Context) {
		c.String(200, "OK")
	})

	router.GET(SERVICE_NAME+"/post", func(c *gin.Context) {

		span := tracer.StartSpan("get post")

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
		span.Finish()

	})

	router.POST(SERVICE_NAME+"/post", func(c *gin.Context) {

		span := tracer.StartSpan("create post")

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
		verified, _ := strconv.ParseBool(user_data["verified"].(string))
		tags := c.PostForm("tags")
		post_tags := make([]string, 1)
		if tags != "" {
			post_tags = strings.Split(tags, ",")
		}
		new_post := &models.Post{

			Postid:       primitive.NewObjectIDFromTimestamp(time.Now()),
			Uid:          user_data["uid"].(string),
			Username:     user_data["username"].(string),
			Screenname:   user_data["screenname"].(string),
			Avatarurl:    user_data["avatarurl"].(string),
			Verified:     verified,
			Imageurl:     img_url,
			Caption:      post_caption,
			Likecount:    0,
			Private:      false,
			Commentcount: 0,
			Viewcount:    0,
			Created:      time.Now(),
			Tag:          post_tags,
		}

		_, create_error := postdb.Create(new_post)
		if create_error == nil {
			c.String(200, "ok")
			span.Finish()
		} else {
			c.String(503, "error")
			span.Finish()
			panic("failed create post")
		}

	})

	router.DELETE(SERVICE_NAME+"/post", func(c *gin.Context) {

		span := tracer.StartSpan("delete post")

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
		span.Finish()

	})

	router.GET(SERVICE_NAME+"/allpost", func(c *gin.Context) {

		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		span := tracer.StartSpan("get all post", ext.RPCServerOption(spanCtx))

		feed_range, query_exist := c.GetQuery("range")
		if query_exist == false {
			feed_range = "8"
		}
		result, findall_err := postdb.FindAll(feed_range)
		if findall_err != nil {
			panic("error getting all post")
		}

		allid_json, json_err := json.Marshal(result)
		if json_err != nil {
			panic("marshal json fail")
		}
		c.String(200, `{ "results": `+string(allid_json)+`}`)
		span.Finish()

	})

	// Internal post endpoint

	router.POST("internal/post", func(c *gin.Context) {

		span := tracer.StartSpan("internal create post")

		img_url := c.PostForm("img_url")
		post_caption := c.PostForm("post_caption")
		uid := c.PostForm("uid")
		username := c.PostForm("username")
		screenname := c.PostForm("screenname")
		avatarurl := c.PostForm("avatarurl")
		tags := c.PostForm("tags")
		post_tags := make([]string, 1)
		if tags != "" {
			post_tags = strings.Split(tags, ",")
		}

		new_post := &models.Post{

			Postid:       primitive.NewObjectIDFromTimestamp(time.Now()),
			Uid:          uid,
			Username:     username,
			Screenname:   screenname,
			Avatarurl:    avatarurl,
			Verified:     false,
			Imageurl:     img_url,
			Caption:      post_caption,
			Likecount:    0,
			Private:      false,
			Commentcount: 0,
			Viewcount:    0,
			Created:      time.Now(),
			Tag:          post_tags,
		}

		_, create_error := postdb.Create(new_post)
		if create_error == nil {
			c.String(200, "ok")
			span.Finish()
		} else {
			c.String(503, "error")
			span.Finish()
			panic("failed create post")
		}

	})

	return router

}

func main() {

	mongo_layer := helpers.NewMongoDatabase()
	postdb := models.NewPostDatabase(mongo_layer)
	authservice := services.NewUserAuthService()
	router := setupRouter(postdb, authservice)
	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}

}
