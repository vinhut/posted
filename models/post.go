package models

import (
	"fmt"
	"github.com/vinhut/posted/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const tableName = "posts"

type PostDatabase interface {
	Find(string, string, interface{}) error
	FindMulti(string, string) ([]string, error)
	FindAll(string) ([]string, error)
	Create(*Post) (bool, error)
	Update() (bool, error)
	Delete(string) (bool, error)
}

type postDatabase struct {
	db helpers.DatabaseHelper
}

type Post struct {
	Postid       primitive.ObjectID `bson:"_id, omitempty"`
	Uid          string
	Username     string
	Screenname   string
	Avatarurl    string
	Verified     bool
	Imageurl     string
	Caption      string
	Likecount    int
	Private      bool
	Commentcount int
	Viewcount    int
	Created      time.Time
	Tag          []string
}

func PostUser() Post {
	post := Post{}
	return post
}

func NewPostDatabase(db helpers.DatabaseHelper) PostDatabase {
	return &postDatabase{
		db: db,
	}
}

func (postdb *postDatabase) Find(column, value string, result_user interface{}) error {
	err := postdb.db.Query(tableName, column, value, result_user)
	if err != nil {
		fmt.Println("model find error ", err)
		return err
	}

	return nil
}

func (postdb *postDatabase) FindMulti(column, value string) ([]string, error) {

	var result_str []string
	data, result_err := postdb.db.FindMulti(tableName, column, value, Post{})

	results := make([]Post, len(data))

	if result_err != nil {
		return nil, result_err
	}

	for i, d := range data {
		if d == nil {
			fmt.Println("d is nil")
		}
		fmt.Println("d = ", d)
		results[i] = d.(Post)
	}

	for _, post := range results {
		result_str = append(result_str, post.Postid.Hex())
	}

	return result_str, nil
}

func (postdb *postDatabase) FindAll(post_range string) ([]string, error) {

	var result_str []string
	data, result_err := postdb.db.FindAll(tableName, post_range, Post{})

	results := make([]Post, len(data))

	if result_err != nil {
		return nil, result_err
	}

	for i, d := range data {
		if d == nil {
			fmt.Println("d is nil")
		}
		fmt.Println("d = ", d)
		results[i] = d.(Post)
	}

	for _, post := range results {
		result_str = append(result_str, post.Postid.Hex())
	}

	return result_str, nil
}

func (postdb *postDatabase) Create(post *Post) (bool, error) {
	err := postdb.db.Insert(tableName, post)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (postdb *postDatabase) Update() (bool, error) {
	return false, nil
}

func (postdb *postDatabase) Delete(postid string) (bool, error) {

	err := postdb.db.Delete(tableName, postid)
	if err != nil {
		return false, err
	}
	return true, nil
}
