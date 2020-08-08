package helpers

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//"net/http"
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
)

type DatabaseHelper interface {
	Query(string, string, string, interface{}) error
	FindAll(string, string, interface{}) ([]interface{}, error)
	Insert(string, interface{}) error
	Delete(string, string) error
}

type MongoDBHelper struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoDatabase() DatabaseHelper {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	log.Print(os.Getenv("MONGO_URL"))
	if err != nil {
		fmt.Println(err)
	}

	db := client.Database(os.Getenv("MONGO_DATABASE"))
	log.Print(os.Getenv("MONGO_DATABASE"))
	return &MongoDBHelper{
		client: client,
		db:     db,
	}
}

func (mdb *MongoDBHelper) Query(collectionName, key, value string, data interface{}) error {

	collection := mdb.db.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	value_hex, value_err := primitive.ObjectIDFromHex(value)
	if value_err != nil {
		return value_err
	}
	result := collection.FindOne(ctx, bson.M{key: value_hex})
	err := result.Decode(data)
	if err != nil {
		fmt.Println("helper mongodb : ", err)
		return err
	}
	return nil
}

func (mdb *MongoDBHelper) FindAll(collectionName string, limit string, obj interface{}) ([]interface{}, error) {

	collection := mdb.db.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	findOptions := options.Find()
	find_limit, _ := strconv.ParseInt(limit, 10, 64)
	findOptions.SetLimit(find_limit)
	cur, err := collection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		fmt.Println("finding fail ", err)
		return nil, err
	}
	defer cur.Close(ctx)

	var container = make([]interface{}, 0)
	for cur.Next(ctx) {

		model := reflect.New(reflect.TypeOf(obj)).Interface()
		decode_err := cur.Decode(model)
		if decode_err != nil {
			fmt.Println("decode fail ", decode_err)
			return nil, decode_err
		}

		fmt.Println("obj = ", obj)
		fmt.Println("model = ", model)
		md := reflect.ValueOf(model).Elem().Interface()
		fmt.Println("md = ", md)
		container = append(container, md)
		fmt.Println("container = ", container)
	}

	return container, nil

}

func (mdb *MongoDBHelper) Insert(collectionName string, data interface{}) error {
	collection := mdb.db.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	new_user, err := bson.Marshal(data)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, new_user)

	if err != nil {
		fmt.Println("Got a real error:", err.Error())
		return err
	}

	return err
}

func (mdb *MongoDBHelper) Delete(collectionName, postid string) error {

	collection := mdb.db.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	postid_hex, objid_err := primitive.ObjectIDFromHex(postid)
	if objid_err != nil {
		return objid_err
	}

	_, err := collection.DeleteOne(ctx, bson.D{{"_id", postid_hex}})
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
